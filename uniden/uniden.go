package uniden

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smoke7385/smk-uniden-bluetooth/types"
	"github.com/smoke7385/smk-uniden-bluetooth/utils"

	"tinygo.org/x/bluetooth"
)

type UnidenCache struct {
	// Owner
	Uniden *Uniden

	// State
	PreMuteVolume         int
	TimeSynced            bool
	RecievedFirstSettings bool
}

func NewUnidenCache(u *Uniden) UnidenCache {
	return UnidenCache{
		// Owner
		Uniden: u,

		// State initializers
		PreMuteVolume:         1,
		TimeSynced:            false,
		RecievedFirstSettings: false,
	}
}

type ConditionalCallbackEvent struct {
	Condition func(u *Uniden) bool
	Callback  func(u *Uniden) error
	Timeout   time.Duration
	Completed bool
}

func (cce *ConditionalCallbackEvent) Unregister() {
	cce.Condition = nil
	cce.Callback = nil
}

type Uniden struct {
	Model   types.Model `validate:"required"`
	Verbose bool

	// Internal state
	server   *UnidenInterfaceServer
	services []*types.Service
	device   *types.Device
	cache    UnidenCache
	address  string

	// State
	Settings Settings
	Alerts   []RadarEvent
	Status   Status

	// Callbacks
	conditionalCallbacks []*ConditionalCallbackEvent
	onServerClientEvent  func(message string)
	onRadarEvent         func(s []RadarEvent)
	onSettingsChange     func(s Settings)
	onStatusUpdate       func(s Status)
	onDisconnect         func()
	onConnect            func()
}

func NewUniden(model types.Model) *Uniden {
	var uniden = Uniden{Model: model, Verbose: true, Settings: defSettings}
	for i := range uniden.Settings {
		uniden.Settings[i].Settings = &uniden.Settings
		uniden.Settings[i].Uniden = &uniden
		uniden.Settings[i].Model = model
	}

	uniden.cache = NewUnidenCache(&uniden)

	return &uniden
}

func (m *Uniden) StayOpen() {
	select {}
}

func (m *Uniden) Connect(address string) error {
	m.println("Connecting to device:", address, "...")

	// Enable bluetooth interface
	utils.Must("enable BLE stack", adapter.Enable())

	// Scan for devices
	result, err := m.scanForDevice(address)
	if err != nil {
		return err
	}

	// Connect to the found device
	_device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
	device := types.Device{Device: _device}
	utils.Must("connect to device device", err)
	if err != nil {
		return err
	}

	// Discover services
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{})
	utils.Must("discover services", err)
	if len(srvcs) == 0 {
		return errors.New("no services identified")
	}

	// Iterate over discovered services and characteristics
	for _, _service := range srvcs {
		service := types.Service{DeviceService: _service}
		device.AddService(&service)

		characteristics, err := service.DiscoverCharacteristics([]bluetooth.UUID{})
		if err != nil {
			println(err)
		}

		for _, _char := range characteristics {
			characteristic := types.Characteristic{DeviceCharacteristic: _char}
			// println("found characteristic for device ", service.UUID().String(), ": ", characteristic.UUID().String())
			service.AddCharacteristic(&characteristic)
			characteristic.AddCallback(m.handleCharacteristicUpdate)
		}

		m.services = append(m.services, &service)
	}

	// Request initial settings data
	sErr := m.requestDeviceState()
	if sErr != nil {
		println("Error getting device state:", sErr)
	} else {
		println("Device state synced successfully")
	}

	// Syncronize the time
	tErr := m.SyncTime()
	if tErr != nil {
		println("Error syncing time:", tErr)
	} else {
		println("Device time synced successfully")
	}

	m.address = address
	m.device = &device

	return nil
}

// scanForDevice scans for the specified device address and returns the result.
func (m *Uniden) scanForDevice(address string) (bluetooth.ScanResult, error) {
	m.println("Scanning for devices...")
	ch := make(chan bluetooth.ScanResult, 1)

	// Start scanning
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == address {
			m.println("Found Uniden device:", result.Address.String(), result.RSSI, result.LocalName())
			adapter.StopScan()
			ch <- result
			return
		}
	})
	if err != nil {
		return bluetooth.ScanResult{}, err
	}

	// Wait for the scan result
	select {
	case result := <-ch:
		return result, nil
	case <-time.After(10 * time.Second): // Timeout after 10 seconds
		return bluetooth.ScanResult{}, errors.New("scan timeout")
	}
}

func (m *Uniden) Disconnect() {
	m.device.Disconnect()

	if m.onDisconnect != nil {
		(m.onDisconnect)()
	}
}

// Callbacks / Listeners
func (m *Uniden) OnRadarEvent(callback func(s []RadarEvent)) {
	m.onRadarEvent = callback
}

func (m *Uniden) OnStatusUpdate(callback func(s Status)) {
	m.onStatusUpdate = callback
}

func (m *Uniden) OnSettingsChange(callback func(s Settings)) {
	m.onSettingsChange = callback
}

func (m *Uniden) OnConnect(callback func()) {
	m.onConnect = callback
}

func (m *Uniden) OnDisconnect(callback func()) {
	m.onDisconnect = callback
}

func (m *Uniden) OnServerClientEvent(callback func(message string)) {
	m.onServerClientEvent = callback
}

// Adapter
func (m *Uniden) StartServer(port int) (*UnidenInterfaceServer, error) {
	m.println("Starting server...")
	server := NewServer(m, port)
	m.server = server

	// Start the server
	go m.server.start()

	return server, nil
}

// Settings
func (m *Uniden) GetSettings() Settings {
	// TODO: This current retrieves the settings state, I may want to change this to
	// "requestSettingsUpdate" to request the settings from the device.
	return m.Settings
}

func (m *Uniden) UpdateSetting(setting string, valueInt int) error {
	settingObj := m.Settings.getByName(setting)

	if settingObj == nil {
		return fmt.Errorf("Settings [%s] not found", setting)
	}

	command := utils.ConcatenateStrings("BTreqSETC:", strconv.Itoa(settingObj.getDeviceStorageIndex()), "=", strconv.Itoa(valueInt))

	// Write the command
	// m.println("Sending command to device: ", command)
	m.SendArbitraryCommand(command)

	return nil
}

// utils
func (m *Uniden) Mute() error {
	vSetting := m.Settings.getByName("Detector volume")
	if vSetting == nil {
		return errors.New("detector volume setting not found")
	}

	// If the volume is already muted, return
	if vSetting.ValueInt == 0 {
		return nil
	}

	pmvCache := vSetting.ValueInt

	// Attempt to mute the device
	err := vSetting.Update(0)

	if err != nil {
		return err
	}

	// If successful, store the old volume
	m.cache.PreMuteVolume = pmvCache

	return nil
}

func (m *Uniden) Unmute() error {
	vSetting := m.Settings.getByName("Detector volume")
	if vSetting == nil {
		return errors.New("detector volume setting not found")
	}

	// If the volume is already unmuted, return
	if vSetting.ValueInt != 0 {
		return nil
	}

	// This should never happen, but if it does, set the volume to 1
	if m.cache.PreMuteVolume == 0 {
		m.cache.PreMuteVolume = 1
	}

	// Attempt to unmute the device
	err := vSetting.Update(m.cache.PreMuteVolume)
	if err != nil {
		return err
	}

	return nil
}

func (m *Uniden) SyncTime() error {
	timeStr := utils.GetDeviceTimeZoneGMT()
	tSetting := m.Settings.getByName("Time zone")

	if tSetting == nil {
		return errors.New("time zone setting not found")
	}

	err, timeInt := tSetting.GetValueInt(timeStr)
	if err != nil {
		return err
	}

	err = tSetting.Update(timeInt)
	if err != nil {
		println("Error syncing time:", err)
	}

	if !m.cache.TimeSynced {
		m.cache.TimeSynced = true
	}

	return nil
}

func (m *Uniden) requestDeviceState() error {
	// Find the command characteristic
	char, err := m.getChar(types.C.Settings.String())
	if err != nil {
		println("Error finding settings characteristic to device")
	}

	// Get value of the settings characteristic
	var data []byte
	_, err = char.Read(data)

	if err != nil {
		return err
	}

	m.handleSettingsUpdate(data, char)

	return nil
}

func (m *Uniden) RegisterConditionalCallback(
	condition func(u *Uniden) bool,
	callback func(u *Uniden) error,
	timeout time.Duration,
) {
	if m.conditionalCallbacks == nil {
		m.conditionalCallbacks = []*ConditionalCallbackEvent{}
	}

	// Create the callback event
	cce := ConditionalCallbackEvent{
		Condition: condition,
		Callback:  callback,
		Timeout:   timeout,
		Completed: false,
	}

	// Add the callback event to the list
	m.conditionalCallbacks = append(m.conditionalCallbacks, &cce)

	// Start the callback event
	go func() {
		time.Sleep(timeout)
		if !cce.Completed {
			cce.Unregister()
		}
	}()
}

func (m *Uniden) runCallbacks() {
	if m.conditionalCallbacks == nil {
		return
	}

	for _, cce := range m.conditionalCallbacks {
		if cce.Condition(m) {
			cce.Callback(m)
			cce.Completed = true

			// Unregister the callback
			cce.Unregister()
		}
	}
}

// TODO: Figure out how marking works.
func (m *Uniden) mark() {}

// Internal event handlers
func (m *Uniden) handleGenericAttribute(buf []byte, c *types.Characteristic) {}

func (m *Uniden) handleSettingsUpdate(buf []byte, c *types.Characteristic) {
	changed := false

	var changedSettings Settings

	for index, value := range buf {

		setting, err := m.Settings.getByDeviceStorageIndex(index)
		if err != nil {
			continue
		}

		if setting.ValueInt != int(value) {
			changedSettings = append(changedSettings, setting)
			setting.ValueInt = int(value)
			changed = true
		}
	}

	if !m.cache.RecievedFirstSettings {
		m.cache.RecievedFirstSettings = true
	}

	if m.onSettingsChange != nil && changed {
		m.runCallbacks()

		if m.server != nil {
			m.server.handleSettingsUpdate(&changedSettings)
		}

		// Invoke the onSettingsChange callback
		(m.onSettingsChange)(m.Settings)
	}
}

func (m *Uniden) handleStatusUpdate(buf []byte, c *types.Characteristic) {
	bStr := string(buf)
	sections := strings.Split(bStr, "&")

	m.Status = Status{
		Voltage: utils.ParseFloat32(sections[0]),
		// What is sections[1]?
		GPS: parseGPS(sections[2]),
		// What is sections[3]?
		Signal: utils.ParseFloat32(sections[4]),
	}

	if m.onStatusUpdate != nil {
		(m.onStatusUpdate)(m.Status)
	}
}

func (m *Uniden) handleRadarEvent(buf []byte, c *types.Characteristic) {
	bStr := string(buf)
	signalSections := strings.Split(bStr, "&")
	var alerts []RadarEvent = m.Alerts

	// TODO: Determine if the position of the signals matters.
	for index, value := range signalSections {
		// This is an empty signal. I suspect the Uniden can only hold 4 signals at a time.
		if value == "0" {
			if len(alerts) > index {
				alerts[index] = RadarEvent{}
			}
			continue
		}

		sections := strings.Split(value, ",")
		// 1,00,K,5,123,24.1090,0,1
		// ?, ?, band, strength, distance???, frequency, ?, ?

		event := RadarEvent{
			Frequency: utils.ParseFloat32(sections[5]),
			Strength:  utils.ParseInt(sections[3]),
			Band:      types.Band(sections[2]),

			LastUpdate: time.Now(),
		}

		if len(alerts) > index {
			alerts[index] = event
		} else {
			alerts = append(alerts, event)
		}

	}

	m.Alerts = alerts

	if m.onRadarEvent != nil {
		(m.onRadarEvent)(alerts)
	}
}

func (m *Uniden) handleResponse(buf []byte, c *types.Characteristic) {}

func (m *Uniden) handleServerClientEvent(message []byte) {
	if m.onServerClientEvent != nil {
		(m.onServerClientEvent)(string(message))
	}
}

func (m *Uniden) handleCharacteristicUpdate(buf []byte, c *types.Characteristic) {
	// m.println("Got data from char: ", types.CharType(c.UUID().String()))
	switch types.CharType(c.UUID().String()) {
	case types.C.GenericAttribute:
		m.handleGenericAttribute(buf, c)
	case types.C.Settings:
		m.handleSettingsUpdate(buf, c)
	case types.C.Status:
		m.handleStatusUpdate(buf, c)
	case types.C.RadarEvent:
		m.handleRadarEvent(buf, c)
	case types.C.Response:
		m.handleResponse(buf, c)
	default:
		m.println("Recieved data from unhandled characteristic:", c.UUID().String())
	}
}

func (m *Uniden) getChar(characteristic string) (*types.Characteristic, error) {
	for _, s := range m.services {
		for _, c := range s.Characteristics {
			if c.UUID().String() == characteristic {
				return &c, nil
			}
		}
	}

	return nil, errors.New("characteristic not found")
}

func (m *Uniden) SendArbitraryCommand(command string) {
	// Find the command characteristic
	char, err := m.getChar(types.C.Command.String())
	if err != nil {
		println("Error writing to device")
	}

	// Write the command
	// m.println("Sending command to device:", command)
	char.WriteWithoutResponse([]byte(command))
}

func (m *Uniden) println(args ...interface{}) {
	if m.Verbose && len(args) > 0 {
		fmt.Println(args...)
	}
}
