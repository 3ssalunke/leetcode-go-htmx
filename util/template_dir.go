package util

import (
	"log"
	"os"
	"path/filepath"
)

func GetTemplateDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("error occured while getting current directory path", err)
		return "", err
	}
	return filepath.Join(pwd, "views"), nil
}
