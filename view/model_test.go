package main

import (
	"fmt"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	currentState.CityNames = append(currentState.CityNames, "Tucson")
	currentState.CityNames = append(currentState.CityNames, "Kochi")

	if err := writeCityNamesToFile(currentState); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
