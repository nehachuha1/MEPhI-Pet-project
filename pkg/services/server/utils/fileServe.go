package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func serve(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		fmt.Printf("Err1 - %v", err)
		return "", err
	}
	defer src.Close()
	dst, err := os.Create("./data/img/" + file.Filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return file.Filename, nil
}

func ServeFiles(files []*multipart.FileHeader) ([]string, error) {
	returnValues := make([]string, 0)
	for _, rawFile := range files {
		fl, err := serve(rawFile)
		if err != nil {
			return []string{}, err
		}
		returnValues = append(returnValues, fl)
	}
	return returnValues, nil
}
