package main

import (
	"./gfycat"
	//"./imgur"
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

// uploadFile uploads a file (image or video) to imgur via their api
func uploadFile(file *os.File, fileType, clientID string) (string, error) {
	// check if fileType parameter is valid for use
	if fileType != "video" && fileType != "image" {
		return "", fmt.Errorf("Error: invalid fileType")
	}

	// create a new body with a multipart form for the form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// create a fileType field (image/video), add the file to it
	fmt.Println(file.Name())
	part, err := writer.CreateFormFile(fileType, file.Name())
	if err != nil {
		return "", err
	}

	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// create new http client, make POST request to upload endpoint
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.imgur.com/3/upload", body)
	if err != nil {
		return "", err
	}

	// add auth and content-type headers
	request.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", clientID))
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// create a new struct for JSON response
	type responseData struct {
		Data struct {
			Error string `json:"error"`
			Link  string `json:"link"`
		} `json:"data"`
		Success bool `json:"success"`
		Status  int  `json:"status"`
	}

	var data responseData

	// encode JSON data into new instance of the struct
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	// if unsuccessful, return the http error code and the error message
	if data.Success == false {
		return "", fmt.Errorf("%d Error: %s", data.Status, data.Data.Error)
	}

	// successful request, return the link to the newly uploaded file
	return data.Data.Link, nil
}
