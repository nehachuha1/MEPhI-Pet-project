package utils

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mephiMainProject/pkg/services/server/database"
	"mime/multipart"
	"os"
)

func encodeImage(filename string, img image.Image, format string) error {
	file, err := os.Create("./data/img/" + filename + "." + format)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, nil)
	case "png":
		return png.Encode(file, img)
	case "gif":
		return gif.Encode(file, img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func decodeImage(fl *multipart.FileHeader) (image.Image, string, error) {
	file, err := fl.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func resizeImage(img image.Image, maxWidth, maxHeight uint) image.Image {
	return resize.Resize(maxWidth, maxHeight, img, resize.Lanczos3)
}

func serve(file *multipart.FileHeader) (string, error) {
	img, format, err := decodeImage(file)
	if err != nil {
		return "", err
	}
	newImg := resizeImage(img, 450, 450)
	newImageName := database.RandStringRunes(16)
	err = encodeImage(newImageName, newImg, format)
	if err != nil {
		return "", err
	}
	return newImageName + "." + format, nil
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

func DeleteFile(filenames []string) error {
	for _, filename := range filenames {
		err := os.Remove("./data/img/" + filename)
		if err != nil {
			return err
		}
	}
	return nil
}
