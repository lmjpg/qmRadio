package main

import (
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func newRadioPopup(window *MainWindowUi, conf *Config) {
	popup := NewDialogUi()

	popup.buttonBox.OnAccepted(func() {
		name, url := popup.nameInput.Text(), popup.urlInput.Text()
		if len(name) == 0 {
			showError("Name cannot be empty.")
		} else if len(url) == 0 {
			showError("Url cannot be empty.")
		} else {
			err := AddRadio(conf, name, url)
			if err != nil {
				showError("Failed to add radio:\n" + err.Error())
			}
			updateRadios(window, conf)
		}
	})

	popup.nameInput.SetFocus()
	popup.Dialog.Show()
}

func startRadio(radio *Radio, urlc chan string) {
	urlc <- radio.Url
}

func updateRadios(window *MainWindowUi, conf *Config) {
	window.RadioList.SetRowCount(len(conf.Radios))
	window.RadioList.SetColumnCount(2)
	for i, radio := range conf.Radios {
		window.RadioList.SetItem(i, 0, qt.NewQTableWidgetItem2(radio.Name))
		window.RadioList.SetItem(i, 1, qt.NewQTableWidgetItem2(radio.Url))
	}
}

func uiFix(window *MainWindowUi) {
	// Apply properties that miqt-uic cannot handle
	buttons := []*qt.QPushButton{window.addButton, window.pauseButton, window.stopButton, window.previousButton, window.nextButton}
	for _, button := range buttons {
		button.SetMinimumSize2(32, 32)
	}

	window.RadioList.VerticalHeader().SetVisible(false)
	window.RadioList.SetHorizontalHeaderLabels([]string{"Name", "Url"})
}

func showError(err string) {
	fmt.Println(err)
	messageBox := qt.NewQMessageBox2()
	messageBox.SetText(err)
	messageBox.Show()
}

func main() {
	qt.NewQApplication(os.Args)
	window := NewMainWindowUi()
	uiFix(window)

	paused := false
	selected := -1

	pausedc := make(chan bool)
	urlc := make(chan string)
	stopc := make(chan int)
	go player(urlc, pausedc, stopc)

	conf, err := GetConfig()
	if err != nil {
		showError("There was an error while loading your saved configuration:\n" + err.Error())
	}

	window.addButton.OnClicked(func() { newRadioPopup(window, conf) })
	window.pauseButton.OnClicked(func() { paused = !paused; pausedc <- paused })
	window.stopButton.OnClicked(func() { stopc <- 1 })
	window.nextButton.OnClicked(func() { selected = (selected + 1) % len(conf.Radios); startRadio(conf.Radios[selected], urlc) })
	window.previousButton.OnClicked(func() {
		if selected == -1 {
			selected = 0
		} else {
			selected = (selected + len(conf.Radios) - 1) % len(conf.Radios)
		}
		startRadio(conf.Radios[selected], urlc)
	})

	updateRadios(window, conf)
	window.RadioList.OnDoubleClicked(func(index *qt.QModelIndex) {
		paused = false
		selected = index.Row()
		startRadio(conf.Radios[selected], urlc)
	})

	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Stretch)
	window.MainWindow.Show()
	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Interactive)
	qt.QApplication_Exec()
}
