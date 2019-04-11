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

// UploadFile uploads a given image file to quad.pe, returns the link to the image
func UploadFile(file *os.File) (string, error) {
	// create a new body with a multipart form for the form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// add return_json field with value true (not having this or setting value to false will return a 405 error)
	err := writer.WriteField("return_json", "true")
	if err != nil {
		return "", err
	}

	// create an image field, add the file to it using io.Copy
	part, err := writer.CreateFormFile("image", file.Name())
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

	request, err := http.NewRequest("POST", "https://quad.pe/api/upload", body)
	if err != nil {
		return "", err
	}

	// set the Content-Type header (request will fail otherwise with a 400 error)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// create a new struct for JSON response
	type quadData struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
		Errors []struct { // only will return one error, but is given as array
			Description string `json:"title"`
		} `json:"errors"`
	}

	var data quadData

	// decode JSON response data into new instance of the struct
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	// if the length of the Errors array is 0, there is no error
	// else, the array is of length 1 and the error is in the "title" field
	if len(data.Errors) != 0 {
		// return new error with the description
		return "", fmt.Errorf("Error: %s", data.Errors[0].Description)
	}

	// return a link with the ID
	return fmt.Sprintf("https://quad.pe/%s", data.Data.ID), nil
}
