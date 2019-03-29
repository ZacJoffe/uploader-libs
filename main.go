package main

import (
	//"./gfycat"
	//"./imgur"
	//"./cmd/gfycat"
	"./cmd/imgur"
	"fmt"
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
	//videoLink := "https://www.youtube.com/watch?v=pf5xjw13mqw"
	type GfycatClient struct {
		ID     string
		Secret string
	}

	type ImgurClient struct {
		ID     string
		Secret string
	}

	/*
		gfyClient := GfycatClient{
			ID:     "2_OUazaV",
			Secret: "vheyue5783LEuIOmwc0A2svpgnFp8Hz7_g5uHXPoRjnn8GwLZBxGoskHQrK4PlxM",
		}

		token, err := gfycat.GenerateToken(gfyClient.ID, gfyClient.Secret)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(token)

		url, err := gfycat.UploadFile("video.mkv", token, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
	*/

	imgurClient := ImgurClient{
		ID:     "0d297558de98a48",
		Secret: "1f6721805889e41a47e797d0f026cbb8a2914b45",
	}

	url, err := imgur.UploadVideo("video.mkv", imgurClient.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)

	link, err := imgur.UploadImage("photo.jpg", imgurClient.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(link)
}
