package main

import "fmt"

func main2() {

	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}
	//adding cities from cityNames.txt
	if err := loadCityNamesFromFile(&currentState); err != nil {
		fmt.Println("Error loading city names into file:", err)
		return
	}

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

	// Iterate through all the values in the WeatherDataMap
	for cityName, cityData := range currentState.WeatherDataMap {
		fmt.Printf("Weather data for %s: %+v\n", cityName, cityData.Pressure)
	}
	// 	// if err := writeCityNamesToFile(currentState); err != nil {
	// 	// 	fmt.Println("Error:", err)
	// 	// 	return
	// 	// }

	// if err := loadCityNamesFromFile(&currentState); err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// 	// Print the loaded city names
	// 	fmt.Println("CityNames:", currentState.CityNames)

}
