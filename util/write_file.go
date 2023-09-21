package util

import (
	"fmt"
	"os"
)

func WriteFile(progLang string, fileExt string, typedCode string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("%s\\docker_runtimes\\%s\\app.%s", wd, progLang, fileExt)

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := []byte(typedCode)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
