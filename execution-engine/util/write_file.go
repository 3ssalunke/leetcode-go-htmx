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

func getFunctionArgString(testCase string) string {
	argsArray := strings.Split(testCase, "\n")
	return strings.Join(argsArray, ", ")
}

func writeExecutionLines(lang string, functionName string, testCases []string) string {
	switch lang {
	case "javascript":
		executionLine := ""
		for _, testCase := range testCases {
			args := getFunctionArgString(testCase)
			executionLine = executionLine + fmt.Sprintf("\nconsole.log(%s(%s))", functionName, args)
		}
		return executionLine
	case "python":
		executionLine := ""
		for _, testCase := range testCases {
			args := getFunctionArgString(testCase)
			executionLine = executionLine + fmt.Sprintf("\nc=Solution()\nprint(c.%s(%s))", functionName, args)
		}
		return executionLine
	default:
		return ""
	}
}

func WriteCodeInExecutionFile(lang string, typedCode string, functionName string, testCases []string) error {
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

	userCodeWithExecutionLines := typedCode + writeExecutionLines(lang, functionName, testCases)

	data := []byte(userCodeWithExecutionLines)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
