package utils

import (
	"net/url"
	"strings"
)

func BuildFileNameFromURL(fullURLFile string) (string, error) {
	fileUrl, err := url.Parse(fullURLFile)
	if err != nil {
		return "", err
	}

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	return fileName, nil
}
