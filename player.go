package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func startAudio(url string) (*beep.Ctrl, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(format)
	log.Println(streamer)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second))
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return ctrl, nil
}

func player(urlc chan string, pausedc chan bool) {
	url := ""
	ctrl := &beep.Ctrl{Streamer: nil, Paused: false}
	for {
		select {
		case url = <-urlc:
			log.Printf("Playing from %v\n", url)
			var err error
			ctrl, err = startAudio(url)
			if err != nil {
				log.Println(err)
				ctrl = &beep.Ctrl{Streamer: nil, Paused: false}
			}
		case paused := <-pausedc:
			ctrl.Paused = paused
		}
	}
}
