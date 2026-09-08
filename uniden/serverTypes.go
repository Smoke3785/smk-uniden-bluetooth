package uniden

// REQUEST STRUCTS
type RequestUpdateSetting struct {
	DeviceStorageIndex int
	ValueInt int
}

func NewRequestUpdateSetting(data ...any) RequestUpdateSetting {
	dataMap := data[0].(map[string]interface{})

	return RequestUpdateSetting{
		DeviceStorageIndex: int(dataMap["deviceStorageIndex"].(float64)),
		ValueInt: int(dataMap["valueInt"].(float64)),
	}
}