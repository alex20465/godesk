package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func AssertNoError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func PositionArgParser(position string) (int, error) {

	var err error
	var inch int
	var cm int

	if strings.HasSuffix(position, "cm") {
		cm, err = strconv.Atoi(strings.Replace(position, "cm", "", 1))

		if err != nil {
			return 0, err
		}

		return (cm * 100) - OFFSET_POSITION, nil
	}

	if strings.HasSuffix(position, "inch") {
		inch, err = strconv.Atoi(strings.Replace(position, "inch", "", 1))

		if err != nil {
			return 0, err
		}

		return int(float64(inch)*2.54*100.0) - OFFSET_POSITION, nil
	}

	return 0, errors.New("failed to parse int")
}
