package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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

func main() {
	clientID := "2_OUazaV"
	clientSecret := "vheyue5783LEuIOmwc0A2svpgnFp8Hz7_g5uHXPoRjnn8GwLZBxGoskHQrK4PlxM"

	token, err := generateToken(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(token)
}
