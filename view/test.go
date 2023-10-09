package main

import "fmt"

func main() {
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}

	// currentState.CityNames = append(currentState.CityNames, "Tucson")
	// currentState.CityNames = append(currentState.CityNames, "Kochi")

	// if err := writeCityNamesToFile(currentState); err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// if err := loadCityNamesFromFile(&currentState); err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// Print the loaded city names
	fmt.Println("CityNames:", currentState.CityNames)

}
