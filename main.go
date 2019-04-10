package main

import (
	//"./gfycat"
	//"./imgur"
	"./quad"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

	photo, err := os.Open("photo.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer photo.Close()

	link, err := quad.UploadFile(photo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(link)
	/*
		type GfycatClient struct {
			ID     string
			Secret string
		}

		type ImgurClient struct {
			ID     string
			Secret string
		}

		gfyClient := GfycatClient{
			ID:     "2_OUazaV",
			Secret: "vheyue5783LEuIOmwc0A2svpgnFp8Hz7_g5uHXPoRjnn8GwLZBxGoskHQrK4PlxM",
		}

		token, err := gfycat.GenerateToken(gfyClient.ID, gfyClient.Secret)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(token)

		video, err := os.Open("video.mp4")
		if err != nil {
			log.Fatal(err)
		}
		defer video.Close()

		url, err := gfycat.UploadFile(video, token, true) // works!
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
	*/

	/*
		imgurClient := ImgurClient{
			ID:     "0d297558de98a48",
			Secret: "1f6721805889e41a47e797d0f026cbb8a2914b45",
		}

		video, err := os.Open("video.mp4")
		if err != nil {
			log.Fatal(err)
		}
		defer video.Close()

		url, err := uploadFile(video, "video", imgurClient.ID) // broken
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)

		photo, err := os.Open("photo.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer photo.Close()

		link, err := uploadFile(photo, "image", imgurClient.ID) // works!
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(link)
	*/
}

func UploadFile(file *os.File) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	err := writer.WriteField("return_json", "true")
	if err != nil {
		return "", err
	}

	fmt.Println(file.Name())
	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		return "", err
	}
	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://quad.pe/api/upload", body)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type quadData struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
		Errors []struct {
			Title string `json:"title"`
		} `json:"errors"`
	}

	var data quadData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	// handle errors

	return fmt.Sprintf("https://quad.pe/%s", data.Data.ID), nil
}
