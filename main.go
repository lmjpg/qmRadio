package main

import (
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func newRadioPopup(conf *Config) {
	popup := NewDialogUi()
	popup.buttonBox.OnAccepted(func() {
		name, url := popup.nameInput.Text(), popup.urlInput.Text()
		err := AddRadio(conf, name, url)
		if err != nil {
			showError("Failed to add radio:\n" + err.Error())
		}
	})
	popup.Dialog.Show()
}

func pauseClicked() {
	fmt.Println("Clicked pause")
}

func uiFix(window *MainWindowUi) {
	// Apply properties that miqt-uic cannot handle
	buttons := []*qt.QPushButton{window.addButton, window.pauseButton, window.stopButton, window.previousButton, window.nextButton}
	for _, button := range buttons {
		button.SetMinimumSize2(32, 32)
	}
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

	window.addButton.OnClicked(func() { newRadioPopup(conf) })
	window.pauseButton.OnClicked(pauseClicked)

	window.MainWindow.Show()
	qt.QApplication_Exec()
}
