package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
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

func getAudioBuffer(url string) (*bufio.Reader, *bufio.Writer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Create("out.mp3")
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(resp.Body)
	filewriter := bufio.NewWriter(file)
	return reader, filewriter, nil
}

func readAudioBuffer(reader *bufio.Reader, filewriter *bufio.Writer) error {
	size := reader.Size()
	buf := make([]byte, size)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}

	log.Printf("Got %v bytes, len %v\n", n, len(buf))
	_, err = filewriter.Write(buf[:n])
	if err != nil {
		return err
	}

	return nil
}

func player(c chan string) {
	var reader *bufio.Reader
	var filewriter *bufio.Writer
	for {
		url, ok := getUrl(c)
		if ok {
			log.Printf("Playing from %v\n", url)
			var err error
			reader, filewriter, err = getAudioBuffer(url)
			if err != nil {
				log.Println(err)
			}
		}

		if reader != nil && filewriter != nil {
			err := readAudioBuffer(reader, filewriter)
			if err != nil {
				log.Println(err)
				reader, filewriter = nil, nil
			}
		}
	}
}
