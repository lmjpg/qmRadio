package main

import (
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func uiFix(window *MainWindowUi) {
	// Apply properties that miqt-uic cannot handle
	buttons := []*qt.QPushButton{window.addButton, window.pauseButton, window.stopButton, window.previousButton, window.nextButton}
	for _, button := range buttons {
		button.SetMinimumSize2(32, 32)
		button.SetMaximumSize2(32, 32)
	}

	window.RadioList.VerticalHeader().SetVisible(false)
	window.RadioList.SetHorizontalHeaderLabels([]string{"Name", "Url"})

	font := qt.NewQFont6("Sans", 10)
	font.SetBold(true)
	window.playerInfo.SetFont(font)
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

	controller := NewController(window)

	go player(controller)

	window.addButton.OnClicked(controller.newRadioPopup)
	window.pauseButton.OnClicked(controller.togglePause)
	window.stopButton.OnClicked(controller.stop)
	window.nextButton.OnClicked(controller.selectNext)
	window.previousButton.OnClicked(controller.selectPrevious)

	controller.updateRadios()
	window.RadioList.OnDoubleClicked(controller.clickRadio)

	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Stretch)
	window.MainWindow.Show()
	window.RadioList.HorizontalHeader().SetSectionResizeMode(qt.QHeaderView__Interactive)
	qt.QApplication_Exec()
}
