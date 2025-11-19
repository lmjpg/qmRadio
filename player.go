package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
)

func startAudio(url string) (*effects.Volume, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Icy-Metadata", "1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	streamer, format, err := decodeAudio(resp)
	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, 0)
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: 0, Silent: false}
	speaker.Play(volume)

	return volume, nil
}

func decodeAudio(resp *http.Response) (beep.Streamer, *beep.Format, error) {
	contentTypeHeader, ok := resp.Header["Content-Type"]
	var contentType string
	if !ok || len(contentTypeHeader) < 1 {
		log.Println("No Content-Type header or invalid Content-Type header, assuming mp3.")
	} else {
		contentType = contentTypeHeader[0]
	}

	icyMetaIntStr, isIcy := resp.Header["Icy-Metaint"]
	var icyMetaInt int
	if isIcy {
		var err error
		icyMetaInt, err = strconv.Atoi(icyMetaIntStr[0])
		if err != nil {
			icyMetaInt = 0
			isIcy = false
		}
	}

	var streamer beep.Streamer
	var format beep.Format
	var err error
	switch contentType {
	case "application/ogg", "audio/ogg", "video/ogg":
		streamer, format, err = vorbis.Decode(MakeHttpStream(resp.Body, isIcy, icyMetaInt))
	case "audio/mpeg":
		fallthrough
	default:
		streamer, format, err = mp3.Decode(MakeHttpStream(resp.Body, isIcy, icyMetaInt))
	}

	if err != nil {
		return nil, nil, err
	}

	return streamer, &format, err
}

func startStream(url string) *effects.Volume {
	log.Printf("Now playing %v\n", url)
	streamer, err := startAudio(url)
	if err != nil {
		log.Println(err)
	}
	return streamer
}

func stopStream(streamer *effects.Volume) {
	if streamer != nil {
		streamer.Streamer.(beep.StreamSeekCloser).Close()
	}
}

func player(urlc chan string, pausedc chan bool, stopc chan int) {
	var url string
	var streamer *effects.Volume
	paused := false
	for {
		select {
		case url = <-urlc:
			stopStream(streamer)
			streamer = startStream(url)

		case paused = <-pausedc:
			if streamer != nil {
				streamer.Silent = paused
			} else if !paused {
				streamer = startStream(url)
			}

		case <-stopc:
			stopStream(streamer)
		}
	}
}
