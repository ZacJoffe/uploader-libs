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

// quadResponseData struct stores the response data for a call to quad's api
type quadResponseData struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	Errors []struct { // only will return one error, but is given as array
		Description string `json:"title"`
	} `json:"errors"`
}

// generateRandomPassword generates a random string of lowercase chars with a size of length
func generateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(97 + rand.Intn(25)) // use ascii values of lowercase chars
	}
	return string(bytes)
}

// UploadImage uploads a given image file to quad.pe, returns the link to the image
func UploadImage(file *os.File) (string, error) {
	imageID, err := upload(file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://quad.pe/%s", imageID), nil
}

// upload uploads a given image file to quad.pe and returns its image ID
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

	// create new http client, add POST request to upload endpoint
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

	// create new data variable for decoding the response body
	var responseData quadResponseData

	// decode JSON response data into new instance of the struct
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", err
	}

	// if the length of the Errors array is 0, there is no error
	// else, the array is of length 1 and the error is in the "title" field
	if len(responseData.Errors) != 0 {
		// return new error with the description
		return "", fmt.Errorf("Error: %s", responseData.Errors[0].Description)
	}

	// return image ID
	return responseData.Data.ID, nil
}

/*
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
*/

// UploadImages creates a new gallery and uploads the given images to it, returns a link to the gallery
func UploadImages(files []*os.File, galleryName string) (string, error) {
	// if empty string was given, give it a default name of "gallery"
	if galleryName == "" {
		galleryName = "gallery"
	}

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

// gallery returns creates a new gallery and uploads adds the imageIDs to it
func gallery(galleryName string, imageIDs []string) (string, error) {
	// create structs for JSON request body
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
	}

	// create a new random password 10 chars long
	password := generateRandomPassword(10)

	// fill out the request data struct
	payload := requestPayload{
		Data: Data{
			Type: "gallery",
			Attributes: Attributes{
				Gallery: fmt.Sprintf("%s!%s", galleryName, password), // TODO: generate password
				Images:  imageIDs,
			},
		},
	}

	// create a new buffer and encode the data
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	// new http client, add PUT request to the gallery endpoint
	client := &http.Client{}

	request, err := http.NewRequest("PUT", "https://quad.pe/api/gallery", body)
	if err != nil {
		return "", err
	}

	// set the Content-Type header (request will fail otherwise with a 400 error)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// create new data variable for decoding the response body
	var responseData quadResponseData

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", err
	}

	// if the length of the Errors array is 0, there is no error
	// else, the array is of length 1 and the error is in the "title" field
	if len(responseData.Errors) != 0 {
		// return new error with the description
		return "", fmt.Errorf("Error: %s", responseData.Errors[0].Description)
	}

	return fmt.Sprintf("https://quad.pe/gallery/#%s", responseData.Data.ID), nil
}
