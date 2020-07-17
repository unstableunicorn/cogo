package lib

import (
	"fmt"
	"log"
	"strings"
)

// HandleAWSError strips and logs strings good format
func HandleAWSError(action string, err error) {
	extractedError := strings.Split(err.Error(), ":")
	errorString := fmt.Sprintf("Error %v,%v", action, extractedError[1])
	log.Fatal(errorString)
}

// GetAWSError strips and returns the error string
func GetAWSError(action string, err error) string {
	extractedErrorStrings := strings.Split(err.Error(), ":")
	return extractedErrorStrings[0]
}
