package utils

import (
	JSON "encoding/json"
	"fmt"
	"strconv"
	"time"
)

func ValueInArray[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func Must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func ParseJSON(json string) (map[string]interface{}, error) {
	var data map[string]interface{}

	return data, JSON.Unmarshal([]byte(json), &data)
}

func ParseFloat32(str string) float32 {
	p, err := strconv.ParseFloat(str, 32)

	if err != nil {
		println("Failed to parse float32", str)
		return ParseFloat32("0")
		// panic("failed to parse float32: " + err.Error())
	}

	return float32(p)
}

func ParseInt(str string) int {
	p, err := strconv.Atoi(str)

	if err != nil {
		panic("failed to parse int: " + err.Error())
	}

	return p
}

func LogStruct(s interface{}) {
	fmt.Printf("%+v\n", s)
}

func SetInterval(callback func(), interval time.Duration) chan bool {
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				callback()
			}
		}
	}()

	return done
}

// IntRange represents a range of integers.
type IntRange struct {
	Start int
	End   int
}

func NewIntRange(start int, end int) *IntRange {
	return &IntRange{
		Start: start,
		End:   end,
	}
}

// Step returns a slice of integers within the specified range with the given step.
func Step(r *IntRange, step int) []int {
	if step <= 0 {
		return nil // Step should be positive
	}

	var result []int
	for i := r.Start; i <= r.End; i += step {
		result = append(result, i)
	}

	return result
}

func GetDeviceTimeZoneGMT() string {
	_, offset := time.Now().Zone()
	return fmt.Sprintf("GMT%+d", offset/3600)
}

func ConcatenateStrings(strings ...string) string {
	var result string
	for _, str := range strings {
		result += str
	}
	return result
}

func After[T any](duration time.Duration, callback func() T) T {
	time.Sleep(duration)
	return callback()
}

func AfterAsync(duration time.Duration, callback func() any) {
	go After(duration, callback)
}
