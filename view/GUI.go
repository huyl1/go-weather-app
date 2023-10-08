package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	//model that holds the current state of the program
	currentState := CurrentState{
		CityNames:      []string{},
		WeatherDataMap: make(map[string]WeatherData),
	}
	//adding Tucson by default
	currentState.CityNames = append(currentState.CityNames, "Tucson")
	tucsonWeather := GetCityData("Tucson")
	currentState.WeatherDataMap["Tucson"] = tucsonWeather

	myApp := app.New()
	myWindow := myApp.NewWindow("Weather App")
	myWindow.Resize(fyne.NewSize(400, 600))

	// settings
	textColor := color.White
	metric := false
	darkMode := false
	fmt.Printf("units: %v\n", metric)
	fmt.Printf("dark mode: %v\n", darkMode)

	// current view variable
	currentCity := "London"

	// String array for cities (for now)
	cities := []string{"London", "Paris", "New York", "Tokyo", "Moscow"}

	// dropdown menu for cities
	citySelect := widget.NewSelect(cities, func(s string) {
		fmt.Println("Selected", s)
	})
	citySelect.SetSelected("London") // Set default city

	// new city input
	newCityInput := widget.NewEntry()
	newCityInput.SetPlaceHolder("Enter City Name to Add")

	addCityButton := widget.NewButton("Add City", func() {
		// TODO: send city name to lookup API, add to text file or sthing
		cities = append(cities, newCityInput.Text)
		currentCity = newCityInput.Text
		citySelect.Options = cities
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
	})

	// toggle button for dark mode
	darkModeToggle := widget.NewCheck("Dark Mode", func(b bool) {
		if b {
			darkMode = true
		} else {
			darkMode = false
		}
	})

	citySelectContainerHorizontal := container.NewHBox(citySelect, tempUnitsToggle, darkModeToggle)
	citySelectContainer := container.NewVBox(citySelectContainerHorizontal)

	// today weather display (big temperature reading next to weather icon)
	todayTemperatureReading := canvas.NewText(" "+"20"+"°", textColor)
	todayTemperatureReading.TextSize = 100
	todayTemperatureReading.Alignment = fyne.TextAlignCenter
	todayTemperatureDescription := canvas.NewText("Sunny", textColor)
	todayTemperatureDescription.Alignment = fyne.TextAlignCenter
	todayTemperatureDescription.TextStyle.Bold = true
	todayTemperatureDescription.TextSize = 20

	// weather forecast for the next three days
	forecast1 := canvas.NewText(" "+"20"+"° ", textColor)
	forecast1.TextSize = 50
	forecast1.Alignment = fyne.TextAlignCenter
	forecast2 := canvas.NewText(" "+"20"+"° ", textColor)
	forecast2.TextSize = 50
	forecast2.Alignment = fyne.TextAlignCenter
	forecast3 := canvas.NewText(" "+"20"+"° ", textColor)
	forecast3.TextSize = 50
	forecast3.Alignment = fyne.TextAlignCenter

	forecastContainer := container.NewHBox(widget.NewLabel("       "), forecast1, forecast2, forecast3)

	// weather details for today (humidity, wind speed, etc.)
	windspeed := widget.NewLabel("Windspeed: 10mph")
	winddir := widget.NewLabel("Wind Direction: N")
	humidity := widget.NewLabel("Humidity: 50%")
	pressure := widget.NewLabel("Pressure: 1000hPa")
	percip := widget.NewLabel("Percipitation: 0%")
	uv := widget.NewLabel("UV Index: 0")

	weatherDetailsContainer1 := container.NewHBox(windspeed, winddir, humidity)
	weatherDetailsContainer2 := container.NewHBox(pressure, percip, uv)

	// last updated time
	lastUpdated := widget.NewLabel("Last Updated: " + time.Now().Format("15:04:05"))
	lastUpdated.Alignment = fyne.TextAlignCenter

	// Assemble the GUI
	mainGUI := citySelectContainer
	mainGUI.Add(cityInputContainer)
	//add 100px padding
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Today"))))
	mainGUI.Add(todayTemperatureDescription)
	mainGUI.Add(todayTemperatureReading)
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Next 3 Days"))))
	mainGUI.Add(forecastContainer)
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Today's Details"))))
	mainGUI.Add(weatherDetailsContainer1)
	mainGUI.Add(weatherDetailsContainer2)
	mainGUI.Add(lastUpdated)
	myWindow.SetContent(mainGUI)

	// go routine to update the weather every 3 seconds (will be 60s in final version)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			updateToday(todayTemperatureReading, todayTemperatureDescription, metric, currentCity)
			updateForecasts(forecast1, forecast2, forecast3, metric, currentCity)
			updateTodayDetails(windspeed, winddir, humidity, pressure, percip, uv, metric, currentCity)
			lastUpdated.SetText("Last Updated: " + time.Now().Format("15:04:05"))
		}
	}()

	myWindow.ShowAndRun()

}

func updateToday(today *canvas.Text, description *canvas.Text, metric bool, currentCity string) {
	fmt.Print("Updating weather for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	// make random number between 20 and 30 (replace this with API call)
	todayTemp := rand.Intn(10) + 20
	today.Text = " " + fmt.Sprint(todayTemp) + "°"
	// set description to random string (replace this with API call)
	todayDescription := "Sunny " + fmt.Sprint(rand.Intn(10))
	description.Text = todayDescription
	today.Refresh()
	description.Refresh()
}

func updateTodayDetails(windspeed *widget.Label, winddir *widget.Label, humidity *widget.Label, pressure *widget.Label, percip *widget.Label, uv *widget.Label, metric bool, currentCity string) {
	// make random number between 20 and 30 (replace this with API call)
	fmt.Print("Updating weather details for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	windspeed.Text = "Windspeed: " + fmt.Sprint(rand.Intn(10))
	winddir.Text = "Wind Direction: " + fmt.Sprint(rand.Intn(10))
	humidity.Text = "Humidity: " + fmt.Sprint(rand.Intn(10))
	pressure.Text = "Pressure: " + fmt.Sprint(rand.Intn(10))
	percip.Text = "Percipitation: " + fmt.Sprint(rand.Intn(10))
	uv.Text = "UV Index: " + fmt.Sprint(rand.Intn(10))
	windspeed.Refresh()
	winddir.Refresh()
	humidity.Refresh()
	pressure.Refresh()
	percip.Refresh()
	uv.Refresh()
}

func updateForecasts(forecast1 *canvas.Text, forecast2 *canvas.Text, forecast3 *canvas.Text, metric bool, currentCity string) {
	fmt.Print("Updating weather for " + currentCity + "\n")
	fmt.Print("metric: " + fmt.Sprint(metric) + "\n")
	// make random number between 20 and 30 (replace this with API call)
	forecast1Temp := rand.Intn(10) + 20
	forecast2Temp := rand.Intn(10) + 20
	forecast3Temp := rand.Intn(10) + 20
	forecast1.Text = " " + fmt.Sprint(forecast1Temp) + "° "
	forecast2.Text = " " + fmt.Sprint(forecast2Temp) + "° "
	forecast3.Text = " " + fmt.Sprint(forecast3Temp) + "° "
	forecast1.Refresh()
	forecast2.Refresh()
	forecast3.Refresh()
}
