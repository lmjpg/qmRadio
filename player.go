package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func getUrl(c chan string) (string, bool) {
	for {
		select {
		case url := <-c:
			return url, true
		default:
			return "", false
		}
	}
}

func getAudioBuffer(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		return err
	}
	log.Println(format)
	log.Println(streamer)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second))
	speaker.Play(streamer)

	return nil
}

func player(c chan string) {
	for {
		url, ok := getUrl(c)
		if ok {
			log.Printf("Playing from %v\n", url)
			err := getAudioBuffer(url)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
