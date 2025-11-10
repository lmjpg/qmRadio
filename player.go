package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func playAudio(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	file, err := os.Create("out.mp3")
	if err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	filewriter := bufio.NewWriter(file)
	for {
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
	}
}

func player(c chan string) {
	i := 0
	for {
		select {
		case url := <-c:
			fmt.Println(url)
			fmt.Println(i)
			log.Println(playAudio(url))
		default:
			i++
		}
	}
}
