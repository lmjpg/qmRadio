package main

import (
	"log"
	"net/http"

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

	speaker.Init(format.SampleRate, 0)
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return ctrl, nil
}

func decodeAudio(resp *http.Response) (beep.Streamer, *beep.Format, error) {
	contentTypeHeader, ok := resp.Header["Content-Type"]
	var contentType string
	if !ok || len(contentTypeHeader) < 1 {
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
	var url string
	var ctrl *beep.Ctrl
	paused := false
	for {
		select {
		case url = <-urlc:
		case paused = <-pausedc:
			if ctrl != nil {
				log.Printf("%v\n", ctrl.Streamer.(beep.StreamSeekCloser).Len())
				if paused {
					speaker.Suspend()
				} else {
					speaker.Resume()
				}
			}
		}

		if !paused {
			if ctrl != nil {
				ctrl.Streamer.(beep.StreamSeekCloser).Close()
			}
			var err error
			ctrl, err = startAudio(url)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
