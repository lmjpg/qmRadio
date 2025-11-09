package main

import "fmt"

func player(c chan string) {
	for true {
		url := <-c
		fmt.Println(url)
	}
}
