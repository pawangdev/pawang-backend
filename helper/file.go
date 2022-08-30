package helper

import "os"

func CreateFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
