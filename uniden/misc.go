package uniden

import (
	"strings"

	"time"

	"github.com/smoke7385/smk-uniden-bluetooth/types"
	"github.com/smoke7385/smk-uniden-bluetooth/utils"
)

type GPSStatus string

const (
	Disconnected GPSStatus = "Disconnected"
	Connected    GPSStatus = "Connected"
	Unknown      GPSStatus = "Unknown"
)

type GPS struct {
	Altitude float32   `json:"altitude"`
	Heading  string    `json:"heading"`
	Speed    float32   `json:"speed"`
	State    GPSStatus `json:"state"`
}

type Status struct {
	Voltage float32 `json:"voltage"`
	Signal  float32 `json:"signal"`
	GPS     GPS
}

// TODO: Clean up this disaster
func (s *Status) Serialize() string {
	tempStruct := struct {
		Voltage *float32 `json:"voltage"`
		Signal  *float32 `json:"signal"`
		GPS     struct {
			Altitude *float32 `json:"altitude"`
			Heading  string   `json:"heading"`
			Speed    *float32 `json:"speed"`
			State    string   `json:"state"`
		}
	}{
		Voltage: &s.Voltage,
		Signal:  &s.Signal,
		GPS: struct {
			Altitude *float32 `json:"altitude"`
			Heading  string   `json:"heading"`
			Speed    *float32 `json:"speed"`
			State    string   `json:"state"`
		}{
			Altitude: &s.GPS.Altitude,
			Heading:  s.GPS.Heading,
			Speed:    &s.GPS.Speed,
			State:    string(s.GPS.State),
		},
	}
	return utils.LooseMarshal(tempStruct)
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

type RadarEvent struct {
	Band      types.Band
	Frequency float32
	Strength  int

	LastUpdate time.Time

	// TODO: Add the rest of the fields.
}

func (r *RadarEvent) Serialize() string {
	tempStruct := struct {
		Band       string     `json:"band"`
		Frequency  *float32   `json:"frequency"`
		Strength   *int       `json:"strength"`
		LastUpdate *time.Time `json:"lastUpdate"`
	}{
		Band:       string(r.Band),
		LastUpdate: &r.LastUpdate,
		Frequency:  &r.Frequency,
		Strength:   &r.Strength,
	}

	return utils.LooseMarshal(tempStruct)
}

type RadarEvents []RadarEvent

func (r *RadarEvents) Serialize() string {
	var serializedEvents []string

	for _, event := range *r {
		serializedEvents = append(serializedEvents, event.Serialize())
	}

	return "[" + strings.Join(serializedEvents, ",") + "]"
}
