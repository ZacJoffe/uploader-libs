package main

import (
	"fmt"
	"github.com/ZacJoffe/gif-uploader/packages"
	"log"
)

func main() {
	/*
		TODO:
			- fix file upload (with sound)
			- fix ... message
			- add optional time selection
			- add videolink as argument parameter
			- if parameter is not given, prompt for import instead of throwing error
			- put link in clipboard
	*/
	//videoLink := "https://www.youtube.com/watch?v=Pf5xjW13MQw"
	clientID := "2_OUazaV"
	clientSecret := "vheyue5783LEuIOmwc0A2svpgnFp8Hz7_g5uHXPoRjnn8GwLZBxGoskHQrK4PlxM"

	token, err := gfycat.GenerateToken(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)

	url, err := gfycat.UploadFile("video.mkv", token, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)

	/*
		url, err := UploadVideo(videoLink, token, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
	*/
}
