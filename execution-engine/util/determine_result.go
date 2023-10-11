package util

import (
	"log"
	"strings"
)

func DetermineResult(containerOutput string, testCaseResults []string) bool {
	log.Println(len(testCaseResults))
	result := true
	eachCaseOutput := strings.Split(containerOutput, "\n")

	for i := 0; i < len(eachCaseOutput)-1; i++ {
		if strings.ReplaceAll(eachCaseOutput[i], " ", "") != testCaseResults[i] {
			result = false
		}
	}

	return result
}
