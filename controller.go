package main

import (
	"github.com/gopxl/beep/v2/effects"
	qt "github.com/mappu/miqt/qt6"
)

type Controller struct {
	Window      *MainWindowUi
	Conf        *Config
	PauseButton *qt.QPushButton
	Paused      bool
	Selected    int
	Streamer    *effects.Volume
	Url         string
}

func NewController(window *MainWindowUi) *Controller {
	conf, err := GetConfig()
	if err != nil {
		showError("There was an error while loading your saved configuration:\n" + err.Error())
	}
	return &Controller{
		Window:      window,
		Conf:        conf,
		PauseButton: window.pauseButton,
		Paused:      false,
		Selected:    0,
		Streamer:    nil,
		Url:         "",
	}
}

func (c *Controller) newRadioPopup() {
	popup := NewDialogUi()

	popup.buttonBox.OnAccepted(func() {
		name, url := popup.nameInput.Text(), popup.urlInput.Text()
		if len(name) == 0 {
			showError("Name cannot be empty.")
		} else if len(url) == 0 {
			showError("Url cannot be empty.")
		} else {
			err := AddRadio(c.Conf, name, url)
			if err != nil {
				showError("Failed to add radio:\n" + err.Error())
			}
			c.updateRadios()
		}
	})

	popup.nameInput.SetFocus()
	popup.Dialog.Show()
}

func (c *Controller) startRadio() {
	radio := c.Conf.Radios[c.Selected]
	c.Url = radio.Url

	stopStream(c.Streamer)
	c.Streamer = startStream(c, c.Url)
}

func (c *Controller) updateRadios() {
	c.Window.RadioList.SetRowCount(len(c.Conf.Radios))
	c.Window.RadioList.SetColumnCount(2)
	for i, radio := range c.Conf.Radios {
		c.Window.RadioList.SetItem(i, 0, qt.NewQTableWidgetItem2(radio.Name))
		c.Window.RadioList.SetItem(i, 1, qt.NewQTableWidgetItem2(radio.Url))
	}
}

func (c *Controller) setPlayerText(text string) {
	c.Window.playerInfo.SetText(text)
	c.Window.MainWindow.SetWindowTitle(text + " - " + qt.QCoreApplication_Tr("qmRadio"))
}

func (c *Controller) setPause(paused bool) {
	c.Paused = paused

	var iconName string
	if paused {
		iconName = "media-playback-start"
	} else {
		iconName = "media-playback-pause"
	}
	icon := qt.QIcon_FromTheme(iconName)
	c.Window.pauseButton.SetIcon(icon)

	if c.Streamer != nil {
		c.Streamer.Silent = paused
	} else if !paused {
		c.Streamer = startStream(c, c.Url)
	}
}

func (c *Controller) togglePause() {
	c.setPause(!c.Paused)
}

func (c *Controller) stop() {
	stopStream(c.Streamer)
}

func (c *Controller) changeSelection(diff int) {
	c.Selected = (c.Selected + diff) % len(c.Conf.Radios)
	c.startRadio()
}

func (c *Controller) selectNext() {
	c.changeSelection(1)
}

func (c *Controller) selectPrevious() {
	c.changeSelection(-1)
}

func (c *Controller) clickRadio(index *qt.QModelIndex) {
	c.setPause(false)
	c.Selected = index.Row()
	c.startRadio()
}
