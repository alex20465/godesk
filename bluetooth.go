package main

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func GetChar(service bluetooth.DeviceService, uuid bluetooth.UUID) (bluetooth.DeviceCharacteristic, error) {
	var found_char bluetooth.DeviceCharacteristic

	chars, err := service.DiscoverCharacteristics(nil)

	if err != nil {
		return found_char, err
	}

	for b := range chars {
		c_uuid := bluetooth.NewUUID((chars[b].UUID().Bytes()))

		if uuid.Bytes() == c_uuid.Bytes() {
			found_char = chars[b]
		}
	}

	return found_char, nil
}

func GetSupportedServices(device bluetooth.Device) (map[[16]byte]bluetooth.DeviceService, error) {

	compatible_services := map[[16]byte]bluetooth.DeviceService{}

	services, err := device.DiscoverServices(nil)

	if err != nil {
		return compatible_services, err
	}

	for i := range EXPECTED_SERVICES {
		for b := range services {

			service_uuid := bluetooth.NewUUID(services[b].Bytes())

			if EXPECTED_SERVICES[i].Bytes() == service_uuid.Bytes() {
				compatible_services[service_uuid.Bytes()] = services[b]
			}
		}
	}

	if len(compatible_services) != 2 {
		return compatible_services, errors.New("Unsupported device.")
	} else {
		return compatible_services, nil
	}
}

func GetService(device *bluetooth.Device, uuid bluetooth.UUID) (bluetooth.DeviceService, error) {
	var service bluetooth.DeviceService

	services, err := GetSupportedServices(*device)

	if err != nil {
		return service, err
	}

	return services[uuid.Bytes()], nil
}

func GoTo(device *bluetooth.Device, desired_position int) error {

	var current_desk_position DeskPosition

	err := TrackDeskPosition(device, &current_desk_position)

	if err != nil {
		return err
	}

	control_service, err := GetService(device, SERVICE_CONTROL)

	if err != nil {
		return err
	}

	control_char, err := GetChar(control_service, SERVICE_CONTROL_CHAR)

	if err != nil {
		return err
	}

	var action []byte

	if current_desk_position.value > uint16(desired_position) {
		action = ACTION_MOVE_DOWN
	} else if current_desk_position.value < uint16(desired_position) {
		action = ACTION_MOVE_UP
	}

	distance := math.Abs(float64(current_desk_position.value) - float64(desired_position))

	last_execution_desk_position := current_desk_position.value

	go control_char.WriteWithoutResponse(action)

	for distance > 30 {

		if reflect.DeepEqual(action, ACTION_MOVE_DOWN) && current_desk_position.speed > 0 {
			err = errors.New("Resistance detected from below the desk")
			break // resistance from down
		} else if reflect.DeepEqual(action, ACTION_MOVE_UP) && current_desk_position.speed < 0 {
			err = errors.New("Resistance detected from above the desk")
			break // resistance from up
		} else if reflect.DeepEqual(action, ACTION_MOVE_DOWN) && current_desk_position.value <= uint16(desired_position) {
			break // finished
		} else if reflect.DeepEqual(action, ACTION_MOVE_UP) && current_desk_position.value >= uint16(desired_position) {
			break // finished
		}

		last_execution_distance := int(math.Abs(float64(last_execution_desk_position) - float64(current_desk_position.value)))

		if last_execution_distance > 50 {
			last_execution_desk_position = current_desk_position.value
			control_char.WriteWithoutResponse(action)
		}

		distance = math.Abs(float64(current_desk_position.value) - float64(desired_position))

		time.Sleep(time.Millisecond * 100)
	}

	control_char.WriteWithoutResponse(ACTION_STOP)

	return err
}

func TrackDeskPosition(device *bluetooth.Device, current_desk_position *DeskPosition) error {
	var err error

	current_position_service, err := GetService(device, SERVICE_CURRENT_POSITION)
	if err != nil {
		return err
	}

	current_position_char, err := GetChar(current_position_service, SERVICE_CURRENT_POSITION_CHAR)

	if err != nil {
		return err
	}
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go current_position_char.EnableNotifications(func(buf []byte) {

		desk_position, err := ParseDeskPositionBytes(buf)

		if err == nil && (*current_desk_position == DeskPosition{}) {
			*current_desk_position = desk_position
			waiter.Done()
		} else if err == nil {
			*current_desk_position = desk_position
		}
	})

	_, err = current_position_char.Read(nil)

	waiter.Wait()

	return err
}

func ConnectDesk(mac string) (*bluetooth.Device, error) {
	var device *bluetooth.Device

	err := adapter.Enable()

	if err != nil {
		return device, err
	}

	mac_address, err := bluetooth.ParseMAC(mac)

	if err != nil {
		return device, err
	}

	address := bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: mac_address}}

	return adapter.Connect(address, bluetooth.ConnectionParams{})
}
