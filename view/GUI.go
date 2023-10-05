package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Weather App")
	myWindow.Resize(fyne.NewSize(400, 600))

	// String array for cities (for now)
	cities := []string{"London", "Paris", "New York", "Tokyo", "Moscow"}

	// Dropdown menu for cities
	citySelect := widget.NewSelect(cities, func(s string) {
		fmt.Println("Selected", s)
	})

	// add textbox for user to type the name of the city, and resize it
	addCityTextbox := widget.NewEntry()
	addCityTextbox.SetPlaceHolder("Enter City Name")

	citySelect.SetSelected("London") // Set default city
	// Add new city button
	addCityButton := widget.NewButton("Add City", func() {
		newCity := addCityTextbox.Text
		cities = append(cities, newCity)
		citySelect.Options = cities
		citySelect.Refresh()
		addCityTextbox.SetText("Enter City Name")
	})

	citySelectContainerHorizontal := container.NewHBox(citySelect, addCityTextbox, addCityButton)
	citySelectContainer := container.NewVBox(citySelectContainerHorizontal)

	// Assemble the GUI
	mainGUI := container.NewVBox(citySelectContainer)

	myWindow.SetContent(mainGUI)
	myWindow.ShowAndRun()
}
