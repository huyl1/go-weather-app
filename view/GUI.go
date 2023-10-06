package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Weather App")
	myWindow.Resize(fyne.NewSize(400, 600))

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
		citySelect.Options = cities
		citySelect.Refresh()
		//clear the input box
		newCityInput.SetText("Enter City Name to Add")
	})

	cityInputContainer := container.NewBorder(nil, nil, nil, addCityButton, newCityInput)
	cityInputContainer.Resize(fyne.NewSize(400, 50))

	// toggle button for temperature units
	tempUnits := widget.NewCheck("Celsius", func(b bool) {
		if b {
			fmt.Println("Celsius")
		} else {
			fmt.Println("Fahrenheit")
		}
	})

	// toggle button for dark mode
	darkMode := widget.NewCheck("Dark Mode", func(b bool) {
		if b {
			fmt.Println("Dark Mode")
		} else {
			fmt.Println("Light Mode")
		}
	})

	citySelectContainerHorizontal := container.NewHBox(citySelect, tempUnits, darkMode)
	citySelectContainer := container.NewVBox(citySelectContainerHorizontal)

	// today weather display (big temperature reading next to weather icon)
	todayTemperatureReading := canvas.NewText(" "+"20"+"°", color.White)
	todayTemperatureReading.TextSize = 100
	todayTemperatureReading.Alignment = fyne.TextAlignCenter

	// weather forecast for the next three days
	forecast1 := canvas.NewText(" "+"20"+"° ", color.White)
	forecast1.TextSize = 50
	forecast1.Alignment = fyne.TextAlignCenter
	forecast2 := canvas.NewText(" "+"20"+"° ", color.White)
	forecast2.TextSize = 50
	forecast2.Alignment = fyne.TextAlignCenter
	forecast3 := canvas.NewText(" "+"20"+"° ", color.White)
	forecast3.TextSize = 50
	forecast3.Alignment = fyne.TextAlignCenter

	forecastContainer := container.NewHBox(widget.NewLabel("       "), forecast1, forecast2, forecast3)

	// Assemble the GUI
	mainGUI := citySelectContainer
	mainGUI.Add(cityInputContainer)
	//add 100px padding
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Today"))))
	mainGUI.Add(todayTemperatureReading)
	mainGUI.Add(container.NewVBox(container.NewVBox(widget.NewLabel("Next 3 Days"))))
	mainGUI.Add(forecastContainer)
	myWindow.SetContent(mainGUI)
	myWindow.ShowAndRun()
}
