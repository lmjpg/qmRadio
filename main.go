package main

import (
	"os"

	qt "github.com/mappu/miqt/qt"
)

func main() {
	qt.NewQApplication(os.Args)
	window := qt.NewQMainWindow2()
	window.SetWindowTitle("qmRadio")
	window.ShowMaximized()
	qt.QApplication_Exec()
}
