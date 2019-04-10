package quad

import (
	"bytes"
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

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	fmt.Println(resp.Body)

	// change
	return "", nil
}
