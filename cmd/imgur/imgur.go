package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadVideo(fileName, clientID string) (string, error) {
	link, err := uploadFile(fileName, "video", clientID)
	if err != nil {
		return "", err
	}

	if link[len(link)-1] == '.' {
		return link[:len(link)-1], nil
	}

	return link, nil
}

func UploadImage(fileName, clientID string) (string, error) {
	return uploadFile(fileName, "image", clientID)
}

func uploadFile(fileName, fileType, clientID string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileType, filepath.Base(file.Name()))
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
