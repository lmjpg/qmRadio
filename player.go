package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
)

func startAudio(url string) (*beep.Ctrl, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	streamer, format, err := decodeAudio(resp)
	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second))
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return ctrl, nil
}

func decodeAudio(resp *http.Response) (beep.Streamer, *beep.Format, error) {
	contentTypeHeader, ok := resp.Header["Content-Type"]
	var contentType string
	if !ok || len(contentTypeHeader) == 1 {
		log.Println("No Content-Type header or invalid Content-Type header, assuming mp3.")
	} else {
		contentType = contentTypeHeader[0]
	}

	var streamer beep.Streamer
	var format beep.Format
	var err error
	switch contentType {
	case "application/ogg", "audio/ogg", "video/ogg":
		streamer, format, err = vorbis.Decode(resp.Body)
	case "audio/mpeg":
		fallthrough
	default:
		streamer, format, err = mp3.Decode(resp.Body)
	}

	if err != nil {
		return nil, nil, err
	}

	return streamer, &format, err
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
