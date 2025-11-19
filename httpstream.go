package main

import (
	"io"
	"log"
	"regexp"
	"strings"
)

// Gets icecast metadata from an audio stream if present.
type HttpStream struct {
	Reader     io.ReadCloser
	IsIcy      bool
	IcyMetaInt int
	SetText    setText
	Pattern    *regexp.Regexp
	Seen       int
}

func MakeHttpStream(body io.ReadCloser, isIcy bool, icyMetaInt int, setPlayerText setText) *HttpStream {
	metaPattern, _ := regexp.Compile(`(\w+)='(.*?)';`)
	return &HttpStream{Reader: body, IsIcy: isIcy, IcyMetaInt: icyMetaInt, SetText: setPlayerText, Pattern: metaPattern, Seen: 0}
}

func (r *HttpStream) Read(p []byte) (n int, e error) {
	n, err := r.Reader.Read(p)
	if err != nil {
		return n, err
	}

	// get icecast metadata
	if r.Seen+n > r.IcyMetaInt {
		metaStart := r.IcyMetaInt - r.Seen
		if metaStart < n && metaStart >= 0 {
			metaLen := int(p[metaStart]) * 16
			metaEnd := metaStart + metaLen
			if metaEnd < n && metaEnd >= metaStart {
				meta := p[metaStart+1 : metaEnd+1]
				strEnd := strings.IndexByte(string(meta), 0)
				if strEnd != -1 {
					meta = meta[:strEnd]
				}

				if len(meta) > 0 {
					r.ParseMetadata(string(meta))
				}

				// remove the metadata section from p
				copy(p[metaStart:], p[metaEnd+1:])

				// resize n to be accurate
				n -= (metaEnd - metaStart) + 1
			}
		}
	}
	r.Seen = (r.Seen + n) % r.IcyMetaInt

	return n, err
}

func (r *HttpStream) Close() error {
	// make this update paused in main.go at some point
	return r.Reader.Close()
}

func (r *HttpStream) ParseMetadata(metadata string) {
	data := r.Pattern.FindAllStringSubmatch(metadata, -1)
	for i := 0; i < len(data); i++ {
		if data[i][1] == "StreamTitle" {
			streamTitle := data[i][2]
			r.SetText(streamTitle)
			log.Printf("Now playing %v\n", streamTitle)
		}
	}
}
