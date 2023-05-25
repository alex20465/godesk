# Go Desk

**godesk** is a Go-lang package CLI tool designed to control your Linak Desk or IKEA IDÅSEN over Bluetooth. The package connects to a low energy actuator system via Bluetooth and enables the remote control of the desk through the command line or a managed internal interface. Simplify your workspace experience with the flexibility and efficiency of godesk.

## Supported Desks

- Linak Desk 8721 (Module)
- IKEA IDÅSEN

## Features

- Move desk Up / Down

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
Prerequisites

You need to have Go installed in your system. If you do not have it installed, follow the instructions [here](https://golang.org/doc/install) to set up.
Installing

To install godesk, you can use go get:

```shell
go get github.com/alex20465/godesk
```

## Usage

### Command Line Usage

To pair or connect your desk:

```
godesk connect
```


To move your desk up:

```
godesk up --mac "XY:XY:XY:XY"
```

To move your desk down:

```
godesk down --mac "XY:XY:XY:XY"
```

## Author

- Alexandros Fotiadis `<fotiadis@alexandros.blue>`

## License

This project is licensed under the MIT License - see the LICENSE.md file for details

## Acknowledgments

- Linak Desk
- IKEA IDÅSEN