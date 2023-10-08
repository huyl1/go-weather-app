package main

import (
	"fmt"
	model "model"
)

func main() {
	currentState := model.CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]model.WeatherData),
	}

	// Populate WeatherData for each city and add it to the map
	newYorkWeather := model.GetCityData("New York")
	currentState.CityNames = append(currentState.CityNames, "New York")
	currentState.WeatherDataMap["New York"] = newYorkWeather

	losAngelesWeather := model.GetCityData("Los Angeles")
	currentState.CityNames = append(currentState.CityNames, "Los Angeles")
	currentState.WeatherDataMap["Los Angeles"] = losAngelesWeather

	// Add more cities and WeatherData as needed

	// To add a new city name and WeatherData:
	cityName := "Chicago"
	houstonWeather := model.GetCityData("Chicago")
	currentState.CityNames = append(currentState.CityNames, cityName)
	currentState.WeatherDataMap[cityName] = houstonWeather

	count := 0
	for _, cityName := range currentState.CityNames {
		// cityName will contain each city name one by one in each iteration
		fmt.Println("City Name:", cityName)
		fmt.Println(count)

		// Access the WeatherData for the current city using the map
		weatherData := currentState.WeatherDataMap[cityName]

		// Access WeatherData fields for the current city
		fmt.Println("Temperature (Celsius):", weatherData.TempC0)
		fmt.Println("Temperature (Fahrenheit):", weatherData.TempF0)
		fmt.Println("Humidity:", weatherData.Humidity)
		// Add more fields as needed

		// Blank line to separate output for different cities
		fmt.Println()
		count++
	}

}
