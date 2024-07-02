package types

import (
	"tinygo.org/x/bluetooth"
)

// DEVICE
type Device struct {
	bluetooth.Device
	Services []Service
}

func (m *Device) AddService(service *Service) {
	service.AddParentDevice(m)
	m.Services = append(m.Services, *service)
}

// SERVICE
type Service struct {
	bluetooth.DeviceService
	Characteristics []Characteristic
	Device          *Device
}

func (m *Service) AddCharacteristic(charicteristic *Characteristic) {
	charicteristic.AddParentService(m)
	m.Characteristics = append(m.Characteristics, *charicteristic)
}

func (m *Service) AddParentDevice(device *Device) {
	m.Device = device
}

// CHARACTERISTIC
type Characteristic struct {
	bluetooth.DeviceCharacteristic
	Service *Service
}

func (m *Characteristic) AddParentService(service *Service) {
	m.Service = service
}

func (m *Characteristic) AddCallback(callback func(buf []byte, m *Characteristic)) {
	m.EnableNotifications(func(buf []byte) {
		callback(buf, m)
	})
}

// Detector model
type Model string

const (
	R4 Model = "R4"
	R8 Model = "R8"
	R9 Model = "R9"
)

// Radar bands
type Band string

const (
	X     Band = "X"
	K     Band = "K"
	Ka    Band = "Ka"
	Laser Band = "Laser"
	Gatso Band = "Gatso"
)

// Characteristic UUIDs
type CharType string

func (c CharType) String() string {
	return string(c)
}

type CharsType struct {
	RadarEvent CharType
	Settings   CharType
	Status     CharType

	Response CharType
	Command  CharType

	GenericAttribute CharType
}

var Characteristics = CharsType{
	// Data channels
	RadarEvent: "6eb675ab-8bd1-1b9a-7444-621e52ec6823",
	Settings:   "2d86686a-53dc-25b3-0c4a-f0e10c8dee20",
	Status:     "6c290d2e-1c03-aca1-ab48-a9b908bae79e",

	// Command / Response channels
	Response: "5987b4ef-3bfa-76a8-e642-92933c31434f",
	Command:  "2c86686a-53dc-25b3-0c4a-f0e10c8dee20",

	// Generic channels
	GenericAttribute: "0000180a-0000-1000-8000-00805f9b34fb",
}

var C = Characteristics
