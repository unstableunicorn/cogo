package cmd

import (
	"fmt"
	"log"
	"strings"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var clientMetadata []string
var userAttributes []string

func mapClientMetadata(cm []string) map[string]*string {
	// Generate the mapped attributes
	mappedClientMetadata := make(map[string]*string)
	if len(cm) > 0 {
		for _, a := range cm {
			values := strings.Split(a, "=")
			if len(values) != 2 {
				errorString := fmt.Sprintf("Invalid client metadata %v \n", a)
				log.Fatal(errorString)
			}
			mappedClientMetadata[values[0]] = &values[1]
		}
	}
	return mappedClientMetadata
}

func mapUserAttributes(ua []string) []*cognito.AttributeType {
	var mappedUserAttributes []*cognito.AttributeType
	if len(ua) > 0 {
		for _, a := range ua {
			values := strings.Split(a, "=")
			if len(values) != 2 {
				errorString := fmt.Sprintf("Invalid user attributes %v \n", a)
				log.Fatal(errorString)
			}

			// aws expects values for booleans to be lower case or other errors can appear
			switch strings.ToLower(values[1]) {
			case "true", "false":
				values[1] = strings.ToLower(values[1])
			}

			attrib := cognito.AttributeType{
				Name:  &values[0],
				Value: &values[1],
			}
			mappedUserAttributes = append(mappedUserAttributes, &attrib)
		}
	}
	return mappedUserAttributes
}
