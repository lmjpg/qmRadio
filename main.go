package main

import (
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func newRadioPopup(window *MainWindowUi) {
	popup := NewDialogUi()
	popup.buttonBox.OnAccepted(func() { fmt.Println(popup.nameInput.Text() + ", " + popup.urlInput.Text()) })
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

func main() {
	qt.NewQApplication(os.Args)
	window := NewMainWindowUi()
	uiFix(window)

	window.addButton.OnClicked(func() { newRadioPopup(window) })
	window.pauseButton.OnClicked(pauseClicked)

	window.MainWindow.Show()
	qt.QApplication_Exec()
}
