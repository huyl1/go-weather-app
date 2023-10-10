//test file that tests all the code in the model.go file
//coverage: 100%
//cd to directory, go test
package main

import (
	"fmt"
	"reflect"
	"testing"
)

//TestState tests the state struct used, whether it is correctly setup and
//if it correctly retrieves the relevant data.
//getCityData is tested here as well (getWeatherForecast is include in getCityData)
//tested features - WeatherData struct, CurrentState struct, getCityData and related functions
func TestState(t *testing.T) {
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	//add some cities to test with
	currentState.CityNames = append(currentState.CityNames, "Tucson")
	currentState.CityNames = append(currentState.CityNames, "Kochi")

	// Create a buffered channel
	dataChannel := make(chan WeatherData, len(currentState.CityNames))

	// Launch Goroutines to fetch and process data for each city
	for _, cityName := range currentState.CityNames {
		go func(cityName string) {
			cityData := getCityData(cityName)
			dataChannel <- cityData // Send the data to the channel
		}(cityName)
	}

	// Receive the data from the channel for each city
	for range currentState.CityNames {
		cityData := <-dataChannel
		// Add cityData to the WeatherDataMap with cityName as the key
		currentState.WeatherDataMap[cityData.CityName] = cityData
	}
	// Check if CityNames list is empty
	if len(currentState.CityNames) == 0 {
		t.Errorf("CityNames list is empty")
	}

	// Check if any of the weather structs have no data (GoodResponse is false)
	for cityName, weatherData := range currentState.WeatherDataMap {
		if !weatherData.GoodResponse {
			t.Errorf("Bad Response, API call failed, No data available for city: %s", cityName)
		}
	}

	// Check if any of the weather structs have no data (uv will be zero
	//which is impossible because uv indices start at 1)
	for cityName, weatherData := range currentState.WeatherDataMap {
		if weatherData.Uv == 0 {
			t.Errorf("UV 0, API call worked, No data available for city: %s", cityName)
		}
	}

}

//TestWrite tests whether the program can save data to disk correctly
//TestWrite sets up the text file for TestRead to test with as well
//getCityData is tested here as well (getWeatherForecast is include in getCityData)
//tested features - WeatherData struct, CurrentState struct, getCityData and related functions,
//                  writeCityNamesToFile
func TestWrite(t *testing.T) {
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	//add some cities to test with
	currentState.CityNames = append(currentState.CityNames, "Tucson")
	currentState.CityNames = append(currentState.CityNames, "Kochi")

	dataChannel := make(chan WeatherData, len(currentState.CityNames)) // Create a buffered channel

	// Launch Goroutines to fetch and process data for each city
	for _, cityName := range currentState.CityNames {
		go func(cityName string) {
			cityData := getCityData(cityName)
			dataChannel <- cityData // Send the data to the channel
		}(cityName)
	}

	// Receive the data from the channel for each city
	for range currentState.CityNames {
		cityData := <-dataChannel
		// Add cityData to the WeatherDataMap with cityName as the key
		currentState.WeatherDataMap[cityData.CityName] = cityData
	}

	if err := writeCityNamesToFile(currentState); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

//TestRead tests whether the program can read data from disk correctly
//TestRead reads the file that TestWrite just made
//getCityData is tested here as well (getWeatherForecast is include in getCityData)
//tested features - WeatherData struct, CurrentState struct, getCityData and related functions,
//                  loadCityNamesToFile
func TestRead(t *testing.T) {
	testState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	//add some cities to test with
	testState.CityNames = append(testState.CityNames, "Tucson")
	testState.CityNames = append(testState.CityNames, "Kochi")

	readDataChannel := make(chan WeatherData, len(testState.CityNames)) // Create a buffered channel

	// Launch Goroutines to fetch and process data for each city
	for _, cityName := range testState.CityNames {
		go func(cityName string) {
			cityData := getCityData(cityName)
			readDataChannel <- cityData // Send the data to the channel
		}(cityName)
	}

	// Receive the data from the channel for each city
	for range testState.CityNames {
		cityData := <-readDataChannel
		// Add cityData to the WeatherDataMap with cityName as the key
		testState.WeatherDataMap[cityData.CityName] = cityData
	}

	//model that holds the current state of the program -- NOEL
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	//adding cities from cityNames.txt
	if err := loadCityNamesFromFile(&currentState); err != nil {
		fmt.Println("Error loading city names into file:", err)
		return
	}

	readDataChannel = make(chan WeatherData, len(currentState.CityNames)) // Create a buffered channel

	// Launch Goroutines to fetch and process data for each city
	for _, cityName := range currentState.CityNames {
		go func(cityName string) {
			cityData := getCityData(cityName)
			readDataChannel <- cityData // Send the data to the channel
		}(cityName)
	}

	// Receive the data from the channel for each city
	for range currentState.CityNames {
		cityData := <-readDataChannel
		// Add cityData to the WeatherDataMap with cityName as the key
		currentState.WeatherDataMap[cityData.CityName] = cityData
	}

	// Compare CityNames from both CurrentState structs
	if !reflect.DeepEqual(testState.CityNames, currentState.CityNames) {
		t.Errorf("CityNames arrays are different:\nExpected: %v\nActual: %v", testState.CityNames, currentState.CityNames)
	}

}
