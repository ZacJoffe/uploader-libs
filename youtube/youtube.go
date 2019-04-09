package youtube

import (
	"fmt"
	//"net/http"
	"strings"
)

func GetVideoID(url string) (string, error) {
	startingIndex := strings.Index(url, "v=")
	if startingIndex == -1 {
		return "", fmt.Errorf("Error: Invalid URL")
	}

	startingIndex += 2

	ampersandIndex := strings.Index(url, "&")
	if ampersandIndex == -1 {
		return url[startingIndex:len(url)], nil
	}

	return url[startingIndex:ampersandIndex], nil
}
