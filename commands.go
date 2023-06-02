package main

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
	"tinygo.org/x/bluetooth"
)

type Context struct {
}

type ConnectCmd struct{}

func (c *ConnectCmd) Run(ctx *Context) error {
	err := adapter.Enable()

	if err != nil {
		return err
	}

	results := map[string]bluetooth.ScanResult{}

	go adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		results[device.Address.String()] = device
	})

	time.Sleep(time.Second * 2)

	adapter.StopScan()

	options := []string{}

	for mac_address := range results {
		options = append(options, mac_address)
	}

	prompt := &survey.Select{
		Message: "Choose found device",
		Options: options,
		Description: func(value string, index int) string {
			return results[value].LocalName()
		},
	}

	var address = ""

	survey.AskOne(prompt, &address)

	var device *bluetooth.Device

	device, err = adapter.Connect(results[address].Address, bluetooth.ConnectionParams{})

	if err != nil {
		return err
	}
	_, err = GetSupportedServices(*device)

	if err != nil {
		return err
	}

	println("Devices connected successfully, use MAC to interact.")

	device.Disconnect()

	return nil
}

type UpCmd struct {
	MAC string `required help:"The mac address of the bluetooth device."`
}

func (mu *UpCmd) Run(ctx *Context) error {

	device, err := ConnectDesk(mu.MAC)

	if err != nil {
		return err
	}

	service, err := GetService(device, SERVICE_CONTROL)

	if err != nil {
		return err
	}

	move_up_char, err := GetChar(service, SERVICE_CONTROL_CHAR)

	if err != nil {
		return err
	}

	_, err = move_up_char.WriteWithoutResponse(ACTION_MOVE_UP)

	device.Disconnect()

	return err
}

type DownCmd struct {
	MAC string `required help:"The mac address of the bluetooth device."`
}

func (mu *DownCmd) Run(ctx *Context) error {

	device, err := ConnectDesk(mu.MAC)

	if err != nil {
		return err
	}

	service, err := GetService(device, SERVICE_CONTROL)

	if err != nil {
		return err
	}

	move_up_char, err := GetChar(service, SERVICE_CONTROL_CHAR)

	if err != nil {
		return err
	}

	_, err = move_up_char.WriteWithoutResponse(ACTION_MOVE_DOWN)

	device.Disconnect()

	return err
}

type GotoCmd struct {
	MAC      string `required help:"The mac address of the bluetooth device."`
	Position string `arg:"" name:"position" help:"Position in CM or Inc"`
}

func (mu *GotoCmd) Run(ctx *Context) error {

	device, err := ConnectDesk(mu.MAC)

	if err != nil {
		return err
	}

	position, err := PositionArgParser(mu.Position)

	if err != nil {
		return err
	}

	err = GoTo(device, position)

	device.Disconnect()

	if err != nil {
		return err
	}

	return nil
}

type StatusCmd struct {
	MAC string `required help:"The mac address of the bluetooth device."`
}

func (mu *StatusCmd) Run(ctx *Context) error {
	device, err := ConnectDesk(mu.MAC)

	if err != nil {
		return err
	}

	var current_desk_position DeskPosition

	err = TrackDeskPosition(device, &current_desk_position)

	in_cm := (current_desk_position.value + uint16(OFFSET_POSITION)) / 100
	in_inch := float64(current_desk_position.value+uint16(OFFSET_POSITION)) / 100.0 / 2.54

	print(`Current hight (absolute): `)

	print(in_cm, `cm | `)
	print(int(in_inch), `inches`)
	println()

	device.Disconnect()

	if err != nil {
		return err
	}

	return nil
}
