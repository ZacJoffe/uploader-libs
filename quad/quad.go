package quad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type quadData struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	Errors []struct { // only will return one error, but is given as array
		Description string `json:"title"`
	} `json:"errors"`
}

func generateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

// UploadFile uploads a given image file to quad.pe, returns the link to the image
func UploadFile(file *os.File) (string, error) {
	imageID, err := upload(file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://quad.pe/%s", imageID), nil
}

func upload(file *os.File) (string, error) {
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

	/*
		// create a new struct for JSON response
		type quadData struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
			Errors []struct { // only will return one error, but is given as array
				Description string `json:"title"`
			} `json:"errors"`
		}
	*/

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

	// return image ID
	return data.Data.ID, nil
}

func NewGallery(galleryName string) (string, error) {
	return gallery(galleryName, []string{})
}

func GalleryAddImage(galleryName string, file *os.File) (string, error) {
	imageID, err := upload(file)
	if err != nil {
		return "", err
	}

	image := []string{imageID}

	return gallery(galleryName, image)
}

func GalleryAddImages(galleryName string, files []*os.File) (string, error) {
	var images []string

	for _, file := range files {
		imageID, err := upload(file)
		if err != nil {
			return "", err
		}

		images = append(images, imageID)
		//time.Sleep(1 * time.Second)
	}

	return gallery(galleryName, images)
}

func gallery(galleryName string, images []string) (string, error) {
	type Attributes struct {
		Gallery string   `json:"gallery"`
		Images  []string `json:"images"`
	}

	type Data struct {
		Type       string     `json:"type"`
		Attributes Attributes `json:"attributes"`
	}

	type requestPayload struct {
		Data Data `json:"data"`
		/*
			Data struct {
				Type       string `json:"type"`
				Attributes struct {
					Gallery string   `json:"gallery"`
					Images  []string `json:"images"`
				} `json:"attributes"`
			} `json:"data"`
		*/
	}

	password := generateRandomPassword(10)
	fmt.Println(password)
	payload := requestPayload{
		Data: Data{
			Type: "gallery",
			Attributes: Attributes{
				Gallery: fmt.Sprintf("%s!%s", galleryName, password), // TODO: generate password
				Images:  images,
			},
		},
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("PUT", "https://quad.pe/api/gallery", body)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var responseData quadData

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://quad.pe/gallery/#%s", responseData.Data.ID), nil
}
