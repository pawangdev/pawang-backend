package helpers

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadImage(filePath, dirPath string, file multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create folder if not exists
	errMkdir := os.MkdirAll(dirPath, os.ModePerm)
	if errMkdir != nil {
		return err
	}

	// Destination Image
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy Image
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
