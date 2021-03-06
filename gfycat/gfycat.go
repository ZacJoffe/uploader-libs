package gfycat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	//"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

/*
	TODO:
		- file upload sound
		- refactor generate token to be private to package (???)
		- fix ... message

		- add optional time selection
		- add videolink as argument parameter
		- if parameter is not given, prompt for import instead of throwing error
		- put link in clipboard
*/

// Generatetoken generates an OAUTH bearer token for given client ID and secret. The token is used for all future REST calls
func GenerateToken(clientID, clientSecret string) (string, error) {
	// create anonymous struct to encode as the payload for the REST call
	payload := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// encode payload for the call
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	// create new HTTP client, add a POST request to the token endpoint with the json payload
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.gfycat.com/v1/oauth/token", body)
	if err != nil {
		return "", err
	}

	// add header for json payload
	request.Header.Add("Content-Type", "application/json")

	// perform the request, encode response into new struct and return the token
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

// UploadVideo uploads a given youtube video, returns link gfycat link
func UploadVideo(link, token string, audio bool) (string, error) {
	// create anonymous payload struct for REST call to gfycat create endpoint
	payload := struct {
		Url   string `json:"fetchUrl"`
		Title string `json:"title"`
		Audio bool   `json:"keepAudio"`
	}{
		Url:   link,
		Title: "test", // remove later
		Audio: audio,
	}

	// encode payload
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return "", err
	}

	// create HTTP client, with POST request to the upload endpoint and appropriate headers
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.gfycat.com/v1/gfycats", body)
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")
	// authentication header too!
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	// do request, encode results in new struct
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

	// output name for debugging
	fmt.Println(data.Gfyname)

	// call GetGyfcatLink to check status of new upload and get the link when it's finished
	url, err := GetGyfcatLink(data.Gfyname, token, audio)
	if err != nil {
		return "", err
	}

	return url, nil
}

// UploadFile uploads a file to gfycat, and returns the gfylink once it's finished uploading.
// Note that the sound doesn't work (yet!!!), doesn't matter what you put for that value.
func UploadFile(file *os.File, token string, audio bool) (string, error) {
	// generate gfyname via a body-less POST request to the "upload" endpoint used in UploadVideo function
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.gfycat.com/v1/gfycats", nil)
	if err != nil {
		return "", err
	}

	// Add oauth token via authentication header
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// encode gfyname into struct
	type responseData struct {
		Gfyname string `json:"gfyname"`
	}

	var data responseData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	fmt.Println(data.Gfyname)

	/*
		// open the given file
		file, err := os.Open(fileName)
		if err != nil {
			return "", err
		}
		defer file.Close()
	*/

	// create a new body for upcoming REST call
	body := new(bytes.Buffer)

	// create a new multipart writer for the form data
	writer := multipart.NewWriter(body)

	// add key field, with a value of the gfyname
	// NOTE: this must come first, BEFORE the file
	err = writer.WriteField("key", data.Gfyname)
	if err != nil {
		return "", err
	}

	// create file field, add the file
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return "", err
	}

	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// make request
	request, err = http.NewRequest("POST", "https://filedrop.gfycat.com", body)
	if err != nil {
		return "", nil
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err = client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	/*
		respRaw, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(respRaw)
		log.Println(bodyString)
	*/

	url, err := GetGyfcatLink(data.Gfyname, token, audio)
	if err != nil {
		return "", err
	}

	return url, nil
}

/*
// UploadFile uploads a file to gfycat, and returns the gfylink once it's finished uploading.
// Note that the sound doesn't work (yet!!!), doesn't matter what you put for that value.
func UploadFile(fileName, token string, audio bool) (string, error) {
	// generate gfyname via a body-less POST request to the "upload" endpoint used in UploadVideo function
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.gfycat.com/v1/gfycats", nil)
	if err != nil {
		return "", err
	}

	// Add oauth token via authentication header
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// encode gfyname into struct
	type responseData struct {
		Gfyname string `json:"gfyname"`
	}

	var data responseData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	fmt.Println(data.Gfyname)

	// open the given file
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// create a new body for upcoming REST call
	body := new(bytes.Buffer)

	// create a new multipart writer for the form data
	writer := multipart.NewWriter(body)

	// add key field, with a value of the gfyname
	// NOTE: this must come first, BEFORE the file
	err = writer.WriteField("key", data.Gfyname)
	if err != nil {
		return "", err
	}

	// create file field, add the file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// make request
	request, err = http.NewRequest("POST", "https://filedrop.gfycat.com", body)
	if err != nil {
		return "", nil
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err = client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	url, err := GetGyfcatLink(data.Gfyname, token, audio)
	if err != nil {
		return "", err
	}

	return url, nil
}
*/

// GetGyfcatLink checks the status of an upload, and returns the url of the webm when encoding is finished
func GetGyfcatLink(gfyname, token string, audio bool) (string, error) {
	// wait 2 seconds, or else it might not be able to find the gfy and will throw an error
	time.Sleep(2 * time.Second)

	// create HTTP client and GET request to status endpoint
	client := &http.Client{}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.gfycat.com/v1/gfycats/fetch/status/%s", gfyname), nil)
	if err != nil {
		return "", err
	}

	// add authentication header
	request.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	// create new struct for responce payload, and a new instance of it
	type responseData struct {
		Task    string `json:"task"`
		Url     string `json:"webmUrl"` // note that this gets the webmUrl, others are available in the response payload
		Gfyname string `json:"gfyname"` // used for if gfy has audio
	}

	// the Url field will not be populated during the initial calls when the gfy is encoding, and the task field will not be populated after encoding has finished
	var data responseData

	// set up flag for loop, and a counter
	status := "encoding"
	count := 1

	for status == "encoding" {
		// do the reqeuest, decode response into struct
		resp, err := client.Do(request)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", err
		}

		// check if gyf is encoding
		switch data.Task {
		case "encoding":
			// if waited an hour, throw an error
			if count == 3600 {
				// new line
				fmt.Printf("\n")
				return "", fmt.Errorf("Gfycat could not be created!")
			}

			// add periods to end of string based on wait time
			fmt.Printf("\rEncoding")
			for i := 1; i <= count; i++ {
				fmt.Printf(".")
			}

			count++

			// wait a second, then make RESt call again
			time.Sleep(1 * time.Second)
		case "NotFoundo":
			//err = fmt.Errorf("Gif: %s not found!", gfyname)
			return "", fmt.Errorf("Gif: %s not found!", gfyname)
		default:
			// new line if gfy was encoding when this was called
			if count > 1 {
				fmt.Printf("\n")
			}
			// set status flag to escape loop
			status = "Done!"
		}
	}

	if audio {
		// concatenate link based off gyfname to return since JSON response with sound doesn't give one
		return fmt.Sprintf("https://gfycat.com/%s", data.Gfyname), nil
	}

	return data.Url, nil
}
