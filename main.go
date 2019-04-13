package main

import (
	//"./gfycat"
	//"./imgur"
	"./quad"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
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
	galleryName := "daswd"

	/*
		link, err := quad.NewGallery(galleryName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(link)
	*/

	var images []*os.File

	photo, err := os.Open("photo.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer photo.Close()

	/*
		link, err := GalleryAddImage(galleryName, photo)
		if err != nil {
			log.Fatal(err)
		}
	*/

	images = append(images, photo)

	link, err := quad.GalleryAddImages(galleryName, images)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(link)
	/*
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
	*/

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

type quadData struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	Errors []struct { // only will return one error, but is given as array
		Description string `json:"title"`
	} `json:"errors"`
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

	time.Sleep(5 * time.Second)
	// return image ID
	return data.Data.ID, nil
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

	payload := requestPayload{
		Data: Data{
			Type: "gallery",
			Attributes: Attributes{
				Gallery: fmt.Sprintf("%s!%s", galleryName, "adawda"), // TODO: generate password
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
