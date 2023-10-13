package util

import (
	"strings"
)

func DetermineResult(containerOutput string, testCaseResults []string) bool {
	result := true
	eachCaseOutput := strings.Split(containerOutput, "\n")

	for i := 0; i < len(eachCaseOutput)-1; i++ {
		if strings.ReplaceAll(eachCaseOutput[i], " ", "") != testCaseResults[i] {
			result = false
		}
	}

	return result
}
