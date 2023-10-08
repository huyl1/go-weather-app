package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WeatherData struct {
	GoodResponse bool
	CityName     string
	TempC        float64 // in Celsius
	TempF        float64
	Humidity     float64 // percentage
	WindMPH      float64
	WindKPH      float64
	PrecipInches float64
	PrecipMm     float64
	Uv           float64
	WindDir      string
	Condition    string
}

func main() {
	cityData := getCityData("New York")
	// Print the data from the CityData struct
	fmt.Println("GoodResponse:", cityData.GoodResponse)
	fmt.Println("CityName:", cityData.CityName)
	fmt.Println("TempC (Celsius):", cityData.TempC)
	fmt.Println("TempF (Fahrenheit):", cityData.TempF)
	fmt.Println("Humidity (%):", cityData.Humidity)
	fmt.Println("WindMPH:", cityData.WindMPH)
	fmt.Println("WindKPH:", cityData.WindKPH)
	fmt.Println("PrecipInches:", cityData.PrecipInches)
	fmt.Println("PrecipMm:", cityData.PrecipMm)
	fmt.Println("Uv:", cityData.Uv)
	fmt.Println("WindDir:", cityData.WindDir)
	fmt.Println("Condition:", cityData.Condition)

}

func getCityData(cityName string) WeatherData {
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

	collectedData.CityName = cityName
	collectedData.TempC = tempC
	collectedData.TempF = tempF
	collectedData.Humidity = humidity
	collectedData.WindMPH = windMPH
	collectedData.WindKPH = windKPH
	collectedData.PrecipInches = precipInches
	collectedData.PrecipMm = precipMm
	collectedData.Uv = uv
	collectedData.Condition = condition
	collectedData.WindDir = windDir
	collectedData.GoodResponse = true

	return collectedData

}
