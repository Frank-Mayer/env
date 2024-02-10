package internal

import (
	"fmt"
	"strings"
)

func UserBool(question string) bool {
	fmt.Printf("%s [y/N]: ", question)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	return response == "y" || response == "Y"
}

func UserChoise(question string, options ...string) int {
	fmt.Printf("%s\n", question)
	for i, option := range options {
		fmt.Printf("[%d] %s\n", i+1, option)
	}
	var response int
	_, err := fmt.Scanln(&response)
	if err != nil {
		return -1
	}
	return response
}

func UserString(question string) string {
	fmt.Printf("%s: ", question)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(response)
}
