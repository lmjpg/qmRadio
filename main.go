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
		err := AddRadio(conf, name, url)
		if err != nil {
			showError("Failed to add radio:\n" + err.Error())
		}
		updateRadios(window, conf)
	})
	popup.nameInput.SetFocus()
	popup.Dialog.Show()
}

func pauseClicked() {
	fmt.Println("Clicked pause")
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

	conf, err := GetConfig()
	if err != nil {
		showError("There was an error while loading your saved configuration:\n" + err.Error())
	}

	window.addButton.OnClicked(func() { newRadioPopup(window, conf) })
	window.pauseButton.OnClicked(pauseClicked)

	updateRadios(window, conf)
	window.RadioList.OnDoubleClicked(func(index *qt.QModelIndex) { fmt.Println(index.Column(), index.Row()) })
	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Stretch)

	window.MainWindow.Show()
	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Interactive)
	qt.QApplication_Exec()
}
