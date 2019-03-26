package main

import (
	//"./gfycat"
	//"./imgur"
	//"./cmd/gfycat"
	"./cmd/imgur"
	"fmt"
	"log"

	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

		url, err := gfycat.UploadFile("video.mkv", token, false)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)

	*/

	imgurClient := ImgurClient{
		ID:     "0d297558de98a48",
		Secret: "1f6721805889e41a47e797d0f026cbb8a2914b45",
	}

	/*
		url, err := UploadFile("video.mkv", imgurClient.ID)
		if err != nil {
			log.Fatal(err)
		}

		if url[len(url)-1] == '.' {
			url = url[:len(url)-1]
		}

		fmt.Println(url)
	*/

	link, err := imgur.UploadImage("photo.jpg", imgurClient.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(link)

}

func UploadFile(fileName, clientID string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}

	io.Copy(part, file)

	/*
		err = writer.WriteField("file", data.Gfyname)
		if err != nil {
			return "", err
		}
	*/

	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.imgur.com/3/upload", body)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", clientID))
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type responseData struct {
		Data struct {
			Error string `json:"error"`
			Link  string `json:"link"`
		} `json:"data"`
		Success bool `json:"success"`
		Status  int  `json:"status"`
	}

	var data responseData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	if data.Success == false {
		return "", fmt.Errorf("%d Error: %s", data.Status, data.Data.Error)
	}

	return data.Data.Link, nil
}
