package main

func main() {
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	// Populate WeatherData for each city and add it to the map

	losAngelesWeather := GetCityData("Istanbul")
	currentState.CityNames = append(currentState.CityNames, "Los Angeles")
	currentState.WeatherDataMap["Los Angeles"] = losAngelesWeather
	print("icon text 0 ", losAngelesWeather.Icon0, "\n")
	print("icon text 1 ", losAngelesWeather.Icon1, "\n")
	print("icon text 2 ", losAngelesWeather.Icon2, "\n")
	print("icon text 3 ", losAngelesWeather.Icon3, "\n")

}
