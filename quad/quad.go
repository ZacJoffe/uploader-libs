package quad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadFile(file *os.File) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	err := writer.WriteField("return_json", "true")
	if err != nil {
		return "", err
	}

	//fmt.Println(file.Name())
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
