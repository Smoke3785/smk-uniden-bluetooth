package uniden

import (
	"encoding/json"
	"errors"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smoke7385/smk-uniden-bluetooth/types"
	"github.com/smoke7385/smk-uniden-bluetooth/utils"
	"tinygo.org/x/bluetooth"
)

// Settings Generators
func generateSpeedValues(
	mphStart int,
	mphEnd int,
	mphSteps int,
	kphStart int,
	kphEnd int,
	kphSteps int,
	withUnit bool,
	valueIsID bool,
) (*Values, *Values) {
	var mph Values
	var kph Values

	suffix := ""

	// MPH Values
	for idx, value := range utils.Step(utils.NewIntRange(mphStart, mphEnd), mphSteps) {
		if withUnit {
			suffix = "mph"
		} else {
			suffix = ""
		}

		key := idx
		if valueIsID {
			key = value
		}

		mph = append(mph, Value{
			Name: strconv.Itoa(value) + suffix,
			ID:   key,
		})
	}

	for idx, value := range utils.Step(utils.NewIntRange(kphStart, kphEnd), kphSteps) {
		if withUnit {
			suffix = "kph"
		} else {
			suffix = ""
		}

		key := idx
		if valueIsID {
			key = value
		}

		kph = append(kph, Value{
			Name: strconv.Itoa(value) + suffix,
			ID:   key,
		})
	}

	return &mph, &kph
}

func generateSlidersRange(min int, max int, step int, suffix string) Values {
	values := Values{}

	for idx, value := range utils.Step(utils.NewIntRange(min, max), step) {
		values = append(values, Value{
			Name: strconv.Itoa(value) + suffix,
			ID:   idx,
		})
	}

	return values
}

func getSpeedValuesLiteral(s *Settings,
	mphStart int,
	mphEnd int,
	mphSteps int,
	kphStart int,
	kphEnd int,
	kphSteps int) *Values {
	mph, kph := generateSpeedValues(mphStart, mphEnd, mphSteps, kphStart, kphEnd, kphSteps, true, true)

	if s.getByName("Speed Units").CurrentValue().Name == "MPH" {
		return mph
	} else {
		return kph
	}

}

func getSpeedValues(s *Settings,
	mphStart int,
	mphEnd int,
	mphSteps int,
	kphStart int,
	kphEnd int,
	kphSteps int) *Values {
	mph, kph := generateSpeedValues(mphStart, mphEnd, mphSteps, kphStart, kphEnd, kphSteps, true, false)

	if s.getByName("Speed Units").CurrentValue().Name == "MPH" {
		return mph
	} else {
		return kph
	}
}

// TODO: Make less retarded - make name first param
func generateBool(indexes ...int) func(name string) *Setting {
	si := map[types.Model]int{}
	setting := Setting{
		Values: BooleanValues,
	}

	for i, v := range indexes {
		if v == -1 {
			continue
		}
		if i == 0 {
			si["R4"] = v
		}
		if i == 1 {
			si["R8"] = v
		}
		if i == 2 {
			si["R9"] = v
		}
	}

	setting.StorageIndex = si

	return func(name string) *Setting {
		setting.Name = name
		return &setting
	}
}

func generateKaSegment(segNum int, r4i int, r8i int, r9i int) *Setting {
	return &Setting{
		Name:   "Ka Segment " + strconv.Itoa(segNum),
		Values: BooleanValues,
		StorageIndex: map[types.Model]int{
			types.R4: (r4i - 1) + segNum,
			types.R8: (r8i - 1) + segNum,
			types.R9: (r9i - 1) + segNum,
		},
	}
}

var adapter = bluetooth.DefaultAdapter

type RadarEvent struct {
	Band      types.Band
	Frequency float32
	Strength  int

	LastUpdate time.Time

	// TODO: Add the rest of the fields.
}

type Settings []*Setting

var BooleanValues = Values{
	Value{"False", 0},
	Value{"True", 1},
}

var ColorValues = Values{
	{"Blue", 0},
	{"Amber", 1},
	{"Green", 2},
	{"Pink", 3},
	{"Gray", 4},
	{"Red", 5},
	{"White", 6},
	{"Purple", 7},
}

var BandColorValues = Values{
	{"Signal strength", 0},
	{"Blue", 1},
	{"Amber", 2},
	{"Green", 3},
	{"Pink", 4},
	{"Gray", 5},
	{"Red", 6},
	{"White", 7},
	{"Purple", 8},
}

var ToneValues = Values{
	Value{"TONE_1", 0},
	Value{"TONE_2", 1},
	Value{"TONE_3", 2},
	Value{"TONE_4", 3},
	Value{"TONE_5", 4},
	Value{"TONE_6", 5},
	Value{"TONE_7", 6},
	Value{"TONE_8", 7},
	Value{"TONE_9", 8},
	Value{"TONE_10", 9},
	Value{"TONE_11", 10},
	Value{"TONE_12", 11},
}

// Turn settings into JSON object to send over the network.
func (s *Settings) Serialize() string {
	var serializedSettings []string

	for _, setting := range *s {
		serializedSettings = append(serializedSettings, setting.Serialize())
	}

	return fmt.Sprintf("[%s]", strings.Join(serializedSettings, ","))
}

func (s *Settings) getByDeviceStorageIndex(index int) (*Setting, error) {
	for _, setting := range *s {
		// setting := &(s)[i]
		if setting.getDeviceStorageIndex() == index {
			return setting, nil
		}
	}

	return nil, errors.New("setting not found")
}

func (s *Settings) getByName(name string) *Setting {
	for _, setting := range *s {
		if strings.EqualFold(setting.Name, name) {
			return setting
		}
	}

	return nil
}

type Value struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"int,omitempty"`
}

func (v *Value) Alt(callback func(v *Value)) {
	callback(v)
}

func (v *Value) Serialize() string {
	return utils.LooseMarshal(v)
}

type Values []Value

func (v *Values) Serialize() string {
	var serializedValues []string

	for _, value := range *v {
		serializedValues = append(serializedValues, value.Serialize())
	}

	return fmt.Sprintf("[%s]", strings.Join(serializedValues, ","))
}

type Setting struct {
	Model        types.Model `validate:"required"`
	Values       Values
	Name         string
	ValueInt     int
	StorageIndex map[types.Model]int
	Settings     *Settings
	DefaultValue any

	Uniden *Uniden

	DynamicValues func(s *Settings) *Values
}

func (s *Setting) Update(valueInt int) error {
	// Validate the value
	err := s.ValidateValueInt(valueInt)
	if err != nil {
		return err
	}

	// Update the value
	err = s.Uniden.UpdateSetting(s.Name, valueInt)
	if err != nil {
		return err
	}

	return nil

}

func (s *Setting) ValidateValueInt(valueInt int) error {
	for _, v := range s.Values {
		if v.ID == valueInt {
			return nil
		}
	}

	return errors.New("value not found")
}

func (s *Setting) GetValueInt(value string) (error, int) {
	for _, v := range *s.GetValues() {
		if strings.EqualFold(v.Name, value) {
			return nil, v.ID
		}
	}

	return errors.New("value not found"), 0
}

func (v *Values) getByInt(id int) *Value {
	for _, value := range *v {
		if value.ID == id {
			return &value
		}
	}

	return &Value{}
}

func (s *Setting) GetValues() *Values {
	if s.DynamicValues != nil {
		return s.DynamicValues(s.Settings)
	}
	return &s.Values
}

func (s *Setting) CurrentValue() *Value {
	return s.GetValues().getByInt(s.ValueInt)
}

func (s *Setting) getDeviceStorageIndex() int {
	// println("Getting storage index for index:", s.Model, s.Name)
	return s.StorageIndex[s.Model]
}

type SerializedSetting struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Values string `json:"values,omitempty"`
}

func (s *Setting) Serialize() string {
	str, err := json.Marshal(SerializedSetting{
		Value:  strconv.Itoa(s.CurrentValue().ID),
		Values: s.Values.Serialize(),
		Name:   s.Name,
	})

	if err != nil {
		return "{}"
	}

	return string(str)
}

// TODO:
// Boolean values don't marshal correctly
