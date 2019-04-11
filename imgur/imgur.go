package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadVideo uploads a given video file to imgur
func UploadVideo(file *os.File, clientID string) (string, error) {
	link, err := uploadFile(file, "video", clientID)
	if err != nil {
		return "", err
	}

	// remove '.' at the end of the returned link (if it exists)
	if link[len(link)-1] == '.' {
		return link[:len(link)-1], nil
	}

	return link, nil
}

// UploadImage uploads a given image file to imgur
func UploadImage(file *os.File, clientID string) (string, error) {
	return uploadFile(file, "image", clientID)
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
	//fmt.Println(file.Name())
	part, err := writer.CreateFormFile(fileType, file.Name())
	if err != nil {
		return "", err
	}
	io.Copy(part, file)

	// close writer explicitly instead of with defer to properly create form
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

	// decode JSON response data into new instance of the struct
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

/*
// uploadFile uploads a file (image or video) to imgur via their api
func uploadFile(fileName, fileType, clientID string) (string, error) {
	// check if fileType parameter is valid for use
	if fileType != "video" && fileType != "image" {
		return "", fmt.Errorf("Error: invalid fileType")
	}

	// open the file
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// create a new body with a multipart form for the form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// create a fileType field (image/video), add the file to it
	part, err := writer.CreateFormFile(fileType, fileName)
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
*/
