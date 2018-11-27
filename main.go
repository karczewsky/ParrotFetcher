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

const (
	parrotDirectory = "./img/"
	workerCount     = 5
)

func worker(jobs chan Emoji, done chan string) {
	fmt.Println("Bringing up worker...")
	for j := range jobs {
		fetchEmote(j, done)
	}
}

func fetchEmote(emoji Emoji, done chan string) {
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

	f, err := os.Create(parrotDirectory + emoji.Name + ".gif")
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
}

func main() {
	parrotData := ParrotsData{}

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

	if _, err := os.Stat(parrotDirectory); os.IsNotExist(err) {
		os.Mkdir(parrotDirectory, 0755)
	}

	jobs := make(chan Emoji, 200)
	done := make(chan string, 200)

	for i := 0; i < workerCount; i++ {
		go worker(jobs, done)
	}

	for _, e := range parrotData.Emojis {
		jobs <- e
	}

	for i := 0; i < len(parrotData.Emojis); i++ {
		println(<-done)
	}

	runServer()
}
