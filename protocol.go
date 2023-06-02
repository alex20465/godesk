package main

import (
	"bytes"
	"encoding/binary"

	"tinygo.org/x/bluetooth"
)

var SERVICE_CONTROL = bluetooth.NewUUID([16]byte{0x8a, 0xf7, 0x15, 0x02, 0x9c, 0x00, 0x49, 0x8a, 0x24, 0x10, 0x8a, 0x33, 0x01, 0x00, 0xfa, 0x99})
var SERVICE_CONTROL_CHAR = bluetooth.NewUUID([16]byte{0x8a, 0xf7, 0x15, 0x02, 0x9c, 0x00, 0x49, 0x8a, 0x24, 0x10, 0x8a, 0x33, 0x02, 0x00, 0xfa, 0x99})
var SERVICE_CURRENT_POSITION = bluetooth.NewUUID([16]byte{0x8a, 0xf7, 0x15, 0x02, 0x9c, 0x00, 0x49, 0x8a, 0x24, 0x10, 0x8a, 0x33, 0x20, 0x00, 0xfa, 0x99})
var SERVICE_CURRENT_POSITION_CHAR = bluetooth.NewUUID([16]byte{0x8a, 0xf7, 0x15, 0x02, 0x9c, 0x00, 0x49, 0x8a, 0x24, 0x10, 0x8a, 0x33, 0x21, 0x00, 0xfa, 0x99})

var ACTION_MOVE_UP = []byte{71, 0}
var ACTION_MOVE_DOWN = []byte{70, 0}
var ACTION_STOP = []byte{255, 0}

var EXPECTED_SERVICES = []bluetooth.UUID{
	SERVICE_CONTROL,
	SERVICE_CURRENT_POSITION,
}

var MAX_POSITION = 6500

var OFFSET_POSITION = 6150

type DeskPosition struct {
	value uint16
	speed int16
}

func ParseDeskPositionBytes(buf []byte) (DeskPosition, error) {
	var desk_position DeskPosition

	reader := bytes.NewBuffer(buf)

	err := binary.Read(reader, binary.LittleEndian, &desk_position.value)

	if err != nil {
		return desk_position, err
	}

	err = binary.Read(reader, binary.LittleEndian, &desk_position.speed)

	if err != nil {
		return desk_position, err
	}

	return desk_position, nil
}
