/*
Golang parrot gif fetching program.
Author: Jakub Karczewski
Args:
	- YAML file containing parrots data struct
Source of YAML: https://cultofthepartyparrot.com/parrotparty.yaml
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ParrotsData struct {
	Title  string
	Emojis []Emoji
}

type Emoji struct {
	Name     string
	Src      string
	Fullname string
}

func main() {
	parrotData := ParrotsData{}
	path := "./img/"
	parrotFile := os.Args[1]

	data, err := ioutil.ReadFile(parrotFile)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &parrotData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 755)
	}

	counter := len(parrotData.Emojis)
	done := make(chan string)

	for _, emoji := range parrotData.Emojis {
		go func(emoji Emoji, done chan string) {
			fmt.Printf("Fetching %v...\n", emoji.Fullname)

			resp, _ := http.Get(emoji.Src)
			f, err := os.Create(path + emoji.Name + ".gif")
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			defer f.Close()

			b := make([]byte, 8)

			for {
				n, err := resp.Body.Read(b)
				f.Write(b[:n])
				if err == io.EOF {
					break
				}
			}
			done <- fmt.Sprintf("Fetched %v", emoji.Fullname)
		}(emoji, done)
	}

	for i := 0; i < counter; i++ {
		println(<-done)
	}
}
