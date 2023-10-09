package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WeatherData struct {
	GoodResponse bool
	CityName     string
	TempC0       float64 // in Celsius
	TempF0       float64
	Humidity     float64 // percentage
	WindMPH      float64
	WindKPH      float64
	PrecipInches float64
	PrecipMm     float64
	Uv           float64
	WindDir      string
	Condition    string
	TempC1       float64
	TempF1       float64 //aaaaaaaa
	TempC2       float64
	TempF2       float64
	TempC3       float64
	TempF3       float64
	Icon0        string
	Icon1        string
	Icon2        string
	Icon3        string
}
type CurrentState struct {
	CityNames      []string               // List of city names as strings
	WeatherDataMap map[string]WeatherData // Map of city names to WeatherData structs
}

// func main() {
// 	dataChannel := make(chan WeatherData)

// 	// Launch a Goroutine to fetch and process the data
// 	go func() {
// 		cityData := getCityData("Tucson")
// 		dataChannel <- cityData // Send the data to the channel
// 	}()

// 	// Receive the data from the channel
// 	cityData := <-dataChannel
// 	// Print the data from the CityData struct
// 	fmt.Println("GoodResponse:", cityData.GoodResponse)

// }

func GetCityData(cityName string) WeatherData {
	var apiKey = "740d078b218647dd88412232230710"
	// Define the API endpoint URL

	encoded := url.QueryEscape(cityName)
	apiUrl := "http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + encoded
	collectedData := WeatherData{GoodResponse: false}

	// Send an HTTP GET request to the API
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return collectedData
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return collectedData
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("API request failed with status code:", response.StatusCode)
		fmt.Println("Response body:", string(responseBody))
		return collectedData
	}

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return collectedData
	}

	tempC := data["current"].(map[string]interface{})["temp_c"].(float64)
	tempF := data["current"].(map[string]interface{})["temp_f"].(float64)
	humidity := data["current"].(map[string]interface{})["humidity"].(float64)
	windMPH := data["current"].(map[string]interface{})["wind_mph"].(float64)
	windKPH := data["current"].(map[string]interface{})["wind_kph"].(float64)
	windDir := data["current"].(map[string]interface{})["wind_dir"].(string)
	precipInches := data["current"].(map[string]interface{})["precip_in"].(float64)
	precipMm := data["current"].(map[string]interface{})["precip_mm"].(float64)
	uv := data["current"].(map[string]interface{})["uv"].(float64)
	condition := data["current"].(map[string]interface{})["condition"].(map[string]interface{})["text"].(string)
	icon0 := data["current"].(map[string]interface{})["condition"].(map[string]interface{})["icon"].(string)

	collectedData.CityName = cityName
	collectedData.TempC0 = tempC
	collectedData.TempF0 = tempF
	collectedData.Humidity = humidity
	collectedData.WindMPH = windMPH
	collectedData.WindKPH = windKPH
	collectedData.PrecipInches = precipInches
	collectedData.PrecipMm = precipMm
	collectedData.Uv = uv
	collectedData.Condition = condition
	collectedData.WindDir = windDir
	collectedData.Icon0 = getImageString(icon0)
	collectedData.GoodResponse = true

	dataChannel := make(chan WeatherData)
	// Launch a Goroutine to fetch and process the data
	go func() {
		getWeatherForecast(apiKey, encoded, &collectedData)
		dataChannel <- collectedData // Send the data to the channel
	}()
	// Receive the data from the channel
	collectedData = <-dataChannel

	return collectedData

}

// Function to fetch weather forecast for the next three days
func getWeatherForecast(apiKey, encoded string, collectedData *WeatherData) {
	apiUrl := "http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + encoded + "&days=3"
	collectedData.GoodResponse = false

	// Send an HTTP GET request to the API
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("API request failed with status code:", response.StatusCode)
		fmt.Println("Response body:", string(responseBody))
	}

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		fmt.Println("Error parsing JSON response:", err)
	}

	// Extract temperature data for the next three days if available
	forecastData, ok := data["forecast"].(map[string]interface{})
	if !ok {
		fmt.Println("No forecast data available.")
	}

	for i := 0; i < 3; i++ {
		dayData, ok := forecastData["forecastday"].([]interface{})[i].(map[string]interface{})["day"].(map[string]interface{})
		if !ok {
			fmt.Println("No data available for day", i+1)
			continue
		}

		avgTempC := dayData["avgtemp_c"].(float64)
		avgTempF := dayData["avgtemp_f"].(float64)
		collectedData.SetTemperature((i + 1), avgTempC, avgTempF)

		iconText := dayData["condition"].(map[string]interface{})["icon"].(string)
		iconText = getImageString(iconText)
		collectedData.SetIcon((i + 1), iconText)
	}

	collectedData.GoodResponse = true
}

func (w *WeatherData) SetTemperature(day int, tempC float64, tempF float64) {
	switch day {
	case 1:
		w.TempC1 = tempC
		w.TempF1 = tempF
	case 2:
		w.TempC2 = tempC
		w.TempF2 = tempF
	case 3:
		w.TempC3 = tempC
		w.TempF3 = tempF
	}
}
func (w *WeatherData) SetIcon(day int, icon string) {
	switch day {
	case 1:
		w.Icon1 = icon
	case 2:
		w.Icon2 = icon
	case 3:
		w.Icon3 = icon
	}
}

func getImageString(path string) string {
	parts := strings.Split(path, "/")
	lastPart := strings.Join(parts[len(parts)-2:], "/")
	return lastPart
}

func loadCityNamesFromFile(currentState *CurrentState) error {
	// Open the file
	filePath := "cityNames.txt"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read each line and add it to the CityNames slice
	for scanner.Scan() {
		cityName := scanner.Text()
		currentState.CityNames = append(currentState.CityNames, cityName)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func writeCityNamesToFile(currentState CurrentState) error {
	filePath := "cityNames.txt"
	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a writer to write to the file
	writer := bufio.NewWriter(file)

	// Write each city name to the file except for the last one
	for i, cityName := range currentState.CityNames {
		_, err := writer.WriteString(cityName)
		if err != nil {
			return err
		}

		// Add a newline character after all but the last city name
		if i < len(currentState.CityNames)-1 {
			_, err := writer.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	// Flush any buffered data to the file
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
