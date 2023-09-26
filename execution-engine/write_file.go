package util

import (
	"fmt"
	"os"
	"strings"
)

func getFileExtension(lang string) string {
	switch lang {
	case "javascript":
		return "js"
	case "python":
		return "py"
	default:
		return ""
	}
}

func writeExecutionLines(lang string, functionName string, testCase string) string {
	args := strings.Split(testCase, "\n")
	argsString := strings.Join(args, ", ")

	switch lang {
	case "javascript":
		return fmt.Sprintf("\n\n\nconsole.log(%s(%s))", functionName, argsString)
	case "python":
		return fmt.Sprintf("\n\n\nc=Solution()\nprint(c.%s(%s))", functionName, argsString)
	default:
		return ""
	}
}

func WriteCodeInExecutionFile(lang string, typedCode string, functionName string, argsString string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("%s\\docker\\runtimes\\%s\\app.%s", wd, lang, getFileExtension(lang))

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	userCodeWithExecutionLines := typedCode + writeExecutionLines(lang, functionName, argsString)

	data := []byte(userCodeWithExecutionLines)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
