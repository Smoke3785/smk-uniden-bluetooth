package uniden

import (
	"strings"

	"github.com/smoke7385/smk-uniden-bluetooth/utils"
)

type GPS struct {
	Altitude float32
	Heading  string
	Speed    float32
	State    string
}

type Status struct {
	Voltage float32
	Signal  float32
	GPS     GPS
}

// Turns the comma-separated GPS data into useful information.
// TODO: Figure out what gpsSections[1] is meant to indicate.
func parseGPS(gpsData string) GPS {
	_GPS := GPS{}
	gpsSections := strings.Split(gpsData, ",")

	_GPS.Heading = gpsSections[0]
	_GPS.Altitude = utils.ParseFloat32(gpsSections[2])

	switch gpsSections[3] {
	case "D":
		_GPS.State = "Disconnected"
	case "C":
		_GPS.State = "Connected"
	default:
		_GPS.State = "Unknown"
	}
	return _GPS
}
