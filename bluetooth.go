package main

import (
	"errors"

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
