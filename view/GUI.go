package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TodayTemp struct {
	today *canvas.Text
	desc  *canvas.Text
}

type Forecast struct {
	forecast1 *canvas.Text
	forecast2 *canvas.Text
	forecast3 *canvas.Text
}

type TodayDetails struct {
	windspeed *widget.Label
	winddir   *widget.Label
	humidity  *widget.Label
	pressure  *widget.Label
	precip    *widget.Label
	uv        *widget.Label
}

func main() {
	//write update functions that update the state vars
	//put in the correct field values at the relevant slots in the GUI
	//on exit, call write function

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

	myApp := app.New()
	myWindow := myApp.NewWindow("Weather App")
	myWindow.Resize(fyne.NewSize(400, 600))

	// settings
	textColor := color.White
	bgcolorLight := color.RGBA{96, 99, 71, 0x80}
	bgColorDark := color.RGBA{96, 99, 71, 0x9A}
	metric := false
	darkMode := true
	myApp.Settings().SetTheme(theme.DarkTheme())
	fmt.Printf("units: %v\n", metric)
	fmt.Printf("dark mode: %v\n", darkMode)

	// current view variable
	currentCity := currentState.CityNames[0]
	currentCityData := currentState.WeatherDataMap[currentCity]

	// String array for cities (for now)
	cities := currentState.CityNames

	// today weather display (big temperature reading next to weather icon)
	// using the Imperial units as default --NOEL
	var tempToday float64
	if metric {
		tempToday = currentCityData.TempC0
	} else {
		tempToday = currentCityData.TempF0
	}
	tempTodayStr := fmt.Sprintf("%.0f", tempToday)
	todayTemperatureReading := canvas.NewText(" "+tempTodayStr+"°", textColor)
	todayTemperatureReading.TextSize = 100
	todayTemperatureReading.Alignment = fyne.TextAlignCenter
	todayTemperatureDescription := canvas.NewText(currentCityData.Condition, textColor)
	todayTemperatureDescription.Alignment = fyne.TextAlignCenter
	todayTemperatureDescription.TextStyle.Bold = true
	todayTemperatureDescription.TextSize = 20

	// weather forecast for the next three days
	var tempDay1 float64
	var tempDay2 float64
	var tempDay3 float64
	if metric {
		tempDay1 = currentCityData.TempC1
		tempDay2 = currentCityData.TempC2
		tempDay3 = currentCityData.TempC3
	} else {
		tempDay1 = currentCityData.TempF1
		tempDay2 = currentCityData.TempF2
		tempDay3 = currentCityData.TempF3
	}
	tempDay1Str := fmt.Sprintf("%.0f", tempDay1)
	tempDay2Str := fmt.Sprintf("%.0f", tempDay2)
	tempDay3Str := fmt.Sprintf("%.0f", tempDay3)
	forecast1 := canvas.NewText(" "+tempDay1Str+"° ", textColor)
	forecast1.TextSize = 50
	forecast1.Alignment = fyne.TextAlignCenter
	forecast2 := canvas.NewText(" "+tempDay2Str+"° ", textColor)
	forecast2.TextSize = 50
	forecast2.Alignment = fyne.TextAlignCenter
	forecast3 := canvas.NewText(" "+tempDay3Str+"° ", textColor)
	forecast3.TextSize = 50
	forecast3.Alignment = fyne.TextAlignCenter

	forecastContainer := container.NewHBox(widget.NewLabel("       "), forecast1, forecast2, forecast3)

	// weather details for today (humidity, wind speed, etc.)
	var windSpeedVal float64
	var windSpeedStr string
	var precipVal float64
	var precipStr string
	if metric {
		windSpeedVal = currentCityData.WindKPH
		windSpeedStr = fmt.Sprintf("%.0fkph", windSpeedVal)
		precipVal = currentCityData.PrecipMm
		precipStr = fmt.Sprintf("%.0fmm", precipVal)
	} else {
		windSpeedVal = currentCityData.WindMPH
		windSpeedStr = fmt.Sprintf("%.0fmph", windSpeedVal)
		precipVal = currentCityData.PrecipInches
		precipStr = fmt.Sprintf("%.0fin", precipVal)
	}

	windspeed := widget.NewLabel("Windspeed: " + windSpeedStr)
	winddir := widget.NewLabel("Wind Direction: " + currentCityData.WindDir)
	humidityStr := fmt.Sprintf("%.0f%", currentCityData.Humidity)
	humidity := widget.NewLabel("Humidity: " + humidityStr)
	pressureStr := fmt.Sprintf("%.0fhPa", currentCityData.Pressure)
	pressure := widget.NewLabel("Pressure: " + pressureStr)
	precip := widget.NewLabel("Precipitation: " + precipStr)
	uvStr := fmt.Sprintf("%.0f", currentCityData.Uv)
	uv := widget.NewLabel("UV Index: " + uvStr)

	weatherDetailsContainerInner := container.NewGridWithRows(3, windspeed, winddir, humidity, pressure, precip, uv)
	weatherDetailsBg := canvas.NewRectangle(bgColorDark)

	weatherDetailsContainerOuter := container.NewStack(weatherDetailsBg, weatherDetailsContainerInner)

	// last updated time
	lastUpdated := widget.NewLabel("Last Updated: " + time.Now().Format("15:04:05"))
	lastUpdated.Alignment = fyne.TextAlignCenter

	// dropdown menu for cities
	citySelect := widget.NewSelect(cities, func(s string) {
		fmt.Println("Selected", s)
		currentCity = s
		currentCityData = currentState.WeatherDataMap[currentCity]
		updateToday(todayTemperatureReading, todayTemperatureDescription, metric, currentCity, currentCityData)
		updateForecasts(forecast1, forecast2, forecast3, metric, currentCity, currentCityData)
		updateTodayDetails(windspeed, winddir, humidity, pressure, precip, uv, metric, currentCity, currentCityData)
	})
	citySelect.SetSelected(currentState.CityNames[0]) // Set default city

	// new city input
	newCityInput := widget.NewEntry()
	newCityInput.SetPlaceHolder("Enter City Name to Add")

	addCityButton := widget.NewButton("Add City", func() {
		// -- NOEL
		currentCity := newCityInput.Text
		currentCityData := getCityData(currentCity)
		if currentCityData.GoodResponse {
			currentState.CityNames = append(currentState.CityNames, currentCity)
			currentState.WeatherDataMap[currentCity] = currentCityData
		} else {
			//TODO: bring up alert saying invalid city name or API may be down
			fmt.Print("Invalid city name provided or the weather API may be down.")
		}
		cities = currentState.CityNames
		citySelect.Options = cities
		citySelect.SetSelected(currentCity)
		citySelect.Refresh()
		//clear the input box
		newCityInput.SetText("Enter City Name to Add")
	})

	cityInputContainer := container.NewBorder(nil, nil, nil, addCityButton, newCityInput)
	cityInputContainer.Resize(fyne.NewSize(400, 50))

	// toggle button for temperature units
	tempUnitsToggle := widget.NewCheck("Metric", func(b bool) {
		if b {
			metric = true
		} else {
			metric = false
		}
		updateToday(todayTemperatureReading, todayTemperatureDescription, metric, currentCity, currentCityData)
		updateForecasts(forecast1, forecast2, forecast3, metric, currentCity, currentCityData)
		updateTodayDetails(windspeed, winddir, humidity, pressure, precip, uv, metric, currentCity, currentCityData)
	})

	// toggle button for dark mode
	darkModeToggle := widget.NewCheck("Dark Mode", func(b bool) {
		if b {
			darkMode = true
			myApp.Settings().SetTheme(theme.DarkTheme())
			textColor = color.White
			todayTemperatureReading.Color = textColor
			todayTemperatureDescription.Color = textColor
			forecast1.Color = textColor
			forecast2.Color = textColor
			forecast3.Color = textColor
			weatherDetailsBg.FillColor = bgColorDark
		} else {
			darkMode = false
			myApp.Settings().SetTheme(theme.LightTheme())
			textColor = color.Black
			todayTemperatureReading.Color = textColor
			todayTemperatureDescription.Color = textColor
			forecast1.Color = textColor
			forecast2.Color = textColor
			forecast3.Color = textColor
			weatherDetailsBg.FillColor = bgcolorLight
		}
	})
	darkModeToggle.SetChecked(true)

	citySelectContainerHorizontal := container.NewHBox(citySelect, tempUnitsToggle, darkModeToggle)
	citySelectContainer := container.NewVBox(citySelectContainerHorizontal)

	// Assemble the GUI
	mainGUI := citySelectContainer
	mainGUI.Add(cityInputContainer)
	//add 100px padding
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Today's Average"))))
	mainGUI.Add(todayTemperatureDescription)
	mainGUI.Add(todayTemperatureReading)
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Next 3 Days' Average"))))
	mainGUI.Add(forecastContainer)
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Real Time Weather Details"))))
	mainGUI.Add(weatherDetailsContainerOuter)
	//mainGUI.Add(weatherDetailsContainer2)
	mainGUI.Add(lastUpdated)
	myWindow.SetContent(mainGUI)

	// go routine to update the weather every 30 seconds (will be 60s in final version)
	go func() {
		for {
			time.Sleep(30 * time.Second)
			//--NOEL
			currentCityData := getCityData(currentCity)
			currentState.WeatherDataMap[currentCity] = currentCityData
			updateToday(todayTemperatureReading, todayTemperatureDescription, metric, currentCity, currentCityData)
			updateForecasts(forecast1, forecast2, forecast3, metric, currentCity, currentCityData)
			updateTodayDetails(windspeed, winddir, humidity, pressure, precip, uv, metric, currentCity, currentCityData)
			lastUpdated.SetText("Last Updated: " + time.Now().Format("15:04:05"))
		}
	}()

	myWindow.ShowAndRun()

}

func updateToday(today *canvas.Text, description *canvas.Text, metric bool, currentCity string, currentCityData WeatherData) {
	fmt.Print("Updating weather for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	var tempToday float64
	if metric {
		tempToday = currentCityData.TempC0
	} else {
		tempToday = currentCityData.TempF0
	}
	tempTodayStr := fmt.Sprintf("%.0f", tempToday)
	today.Text = " " + tempTodayStr + "°"
	// set description to random string (replace this with API call)
	description.Text = currentCityData.Condition
	today.Refresh()
	description.Refresh()
}

func updateTodayDetails(windspeed *widget.Label, winddir *widget.Label, humidity *widget.Label, pressure *widget.Label,
	precip *widget.Label, uv *widget.Label, metric bool, currentCity string, currentCityData WeatherData) {
	// make random number between 20 and 30 (replace this with API call)
	fmt.Print("Updating weather details for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	var windSpeedVal float64
	var windSpeedStr string
	var precipVal float64
	var precipStr string
	if metric {
		windSpeedVal = currentCityData.WindKPH
		windSpeedStr = fmt.Sprintf("%.0fkph", windSpeedVal)
		precipVal = currentCityData.PrecipMm
		precipStr = fmt.Sprintf("%.0fmm", precipVal)
	} else {
		windSpeedVal = currentCityData.WindMPH
		windSpeedStr = fmt.Sprintf("%.0fmph", windSpeedVal)
		precipVal = currentCityData.PrecipInches
		precipStr = fmt.Sprintf("%.0fin", precipVal)
	}
	humidityStr := fmt.Sprintf("%.0f%", currentCityData.Humidity)
	pressureStr := fmt.Sprintf("%.0fhPa", currentCityData.Pressure)
	uvStr := fmt.Sprintf("%.0f", currentCityData.Uv)

	windspeed.Text = "Windspeed: " + windSpeedStr
	winddir.Text = "Wind Direction: " + currentCityData.WindDir
	humidity.Text = "Humidity: " + humidityStr
	pressure.Text = "Pressure: " + pressureStr
	precip.Text = "Percipitation: " + precipStr
	uv.Text = "UV Index: " + uvStr
	windspeed.Refresh()
	winddir.Refresh()
	humidity.Refresh()
	pressure.Refresh()
	precip.Refresh()
	uv.Refresh()
}

func updateForecasts(forecast1 *canvas.Text, forecast2 *canvas.Text, forecast3 *canvas.Text, metric bool,
	currentCity string, currentCityData WeatherData) {
	fmt.Print("Updating weather for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	var tempDay1 float64
	var tempDay2 float64
	var tempDay3 float64
	if metric {
		tempDay1 = currentCityData.TempC1
		tempDay2 = currentCityData.TempC2
		tempDay3 = currentCityData.TempC3
	} else {
		tempDay1 = currentCityData.TempF1
		tempDay2 = currentCityData.TempF2
		tempDay3 = currentCityData.TempF3
	}
	tempDay1Str := fmt.Sprintf("%.0f", tempDay1)
	tempDay2Str := fmt.Sprintf("%.0f", tempDay2)
	tempDay3Str := fmt.Sprintf("%.0f", tempDay3)
	forecast1.Text = " " + tempDay1Str + "° "
	forecast2.Text = " " + tempDay2Str + "° "
	forecast3.Text = " " + tempDay3Str + "° "
	forecast1.Refresh()
	forecast2.Refresh()
	forecast3.Refresh()
}

//--NOEL when closing the program just put this line: writeCityNamesToFile(currentState)
