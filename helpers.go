package main

import "os"

func AssertNoError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
