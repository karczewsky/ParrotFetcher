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
	"time"

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

	if len(os.Args) < 2 {
		log.Fatal("Please point parrot file as command argument")
	}
	parrotFile := os.Args[1]

	data, err := ioutil.ReadFile(parrotFile)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading file %v \nError: %v", parrotFile, err))
	}

	err = yaml.Unmarshal([]byte(data), &parrotData)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error unmarshaling file %v \nError: %v", parrotFile, err))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	counter := len(parrotData.Emojis)
	done := make(chan string)

	for _, emoji := range parrotData.Emojis {
		time.Sleep(40 * time.Millisecond)
		go func(emoji Emoji, done chan string) {
			fmt.Printf("Fetching %v...\n", emoji.Fullname)

			resp, err := http.Get(emoji.Src)
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				done <- fmt.Sprintf("Error fetching %v CODE: %v URL: %v ", emoji.Name, resp.StatusCode, emoji.Src)
				return
			}
			if err != nil {
				done <- fmt.Sprintf("Error fetching %v", emoji.Name)
				return
			}

			f, err := os.Create(path + emoji.Name + ".gif")
			if err != nil {
				done <- fmt.Sprintf("Error creating file for %v", emoji.Name)
				return
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
