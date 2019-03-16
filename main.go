package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func generateToken(clientID, clientSecret string) (string, error) {
	payload := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", "https://api.gfycat.com/v1/oauth/token", body)
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type responseData struct {
		Token string `json:"access_token"`
	}

	var tokenResponse responseData

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Token, nil
}

func uploadVideo(link, token string) (string, error) {
	payload := struct {
		Url   string `json:"fetchUrl"`
		Title string `json:"title"`
	}{
		Url:   link,
		Title: "test",
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.gfycat.com/v1/gfycats", body)
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type responseData struct {
		Gfyname string `json:"gfyname"`
	}

	var data responseData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Gfyname, nil
}

func checkStatus(gfyname, token string) (string, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.gfycat.com/v1/gfycats/fetch/status/%s", gfyname), nil)
	if err != nil {
		return "", err
	}

	//request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	type responseData struct {
		Task string `json:"task"`
		Url  string `json:"webmUrl"`
	}
	var data responseData

	status := "encoding"
	count := 1

	for status == "encoding" {
		resp, err := client.Do(request)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", err
		}

		switch data.Task {
		case "encoding":
			// if waited 30 seconds, throw an error
			if count == 30 {
				// new line
				fmt.Printf("\n")
				return "", fmt.Errorf("Gfycat could not be created!")
			}

			fmt.Printf("\rEncoding")
			for i := 1; i <= count; i++ {
				fmt.Printf(".")
			}
			count++
			//fmt.Printf("\rEncoding...")
			time.Sleep(1 * time.Second)
		case "NotFoundo":
			//err = fmt.Errorf("Gif: %s not found!", gfyname)
			return "", fmt.Errorf("Gif: %s not found!", gfyname)
		default:
			if count > 1 {
				fmt.Printf("\n")
			}
			status = "Done!"
		}
	}

	return data.Url, nil
}

func main() {
	videoLink := "https://www.youtube.com/watch?v=Pf5xjW13MQw"
	clientID := "2_OUazaV"
	clientSecret := "vheyue5783LEuIOmwc0A2svpgnFp8Hz7_g5uHXPoRjnn8GwLZBxGoskHQrK4PlxM"

	token, err := generateToken(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)

	gfyname, err := uploadVideo(videoLink, token)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gfyname)

	url, err := checkStatus(gfyname, token)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)
}
