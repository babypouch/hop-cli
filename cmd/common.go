package cmd

import (
	"net/url"
	"strings"
)

func buildFileNameFromURL(fullURLFile string) (string, error) {
	fileUrl, err := url.Parse(fullURLFile)
	if err != nil {
		return "", err
	}

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	return fileName, nil
}
