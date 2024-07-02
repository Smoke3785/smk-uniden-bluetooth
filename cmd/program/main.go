package main

import (
	"strconv"
	"time"

	"github.com/smoke7385/smk-uniden-bluetooth/uniden"
	"github.com/smoke7385/smk-uniden-bluetooth/utils"
)

func main() {
	unidenInstance := uniden.NewUniden("R4")

	unidenInstance.OnStatusUpdate(func(status uniden.Status) {
		// println("Status updated:")
		// Utils.LogStruct(status)
	})

	unidenInstance.OnRadarEvent(func(events []uniden.RadarEvent) {
		println("Alerts:")
		utils.LogStruct(events)
	})

	unidenInstance.OnSettingsChange(func(settings uniden.Settings) {
		// for _, setting := range settings {
		// 	println("Setting:", setting.Name, "Value:", setting.CurrentValue().Name, "(", setting.CurrentValue().ID, ")")
		// }
	})

	err := unidenInstance.Connect("E0:00:00:00:4F:C5")

	if err != nil {
		println("Failed to connect to device:", err)
		return
	}

	// unidenInstance.SendArbitraryCommand("BTreqSETC:50=2")

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	unidenInstance.Unmute()
	// }()

	_, err = unidenInstance.StartServer(8080)

	go test(unidenInstance)
	unidenInstance.StayOpen()
}

var i int = 0

func ic() int {
	if i >= 7 {
		i = 0
	} else {
		i++
	}
	return i
}

func test(uniden *uniden.Uniden) {
	// Call setInterval with the desired interval and callback function
	utils.SetInterval(func() {
		str := "BTreqSETC:50=" + strconv.Itoa(ic())
		uniden.SendArbitraryCommand(str)
		// println()
	}, 500*time.Millisecond)
}
