package cmd

import (
	"fmt"
	"log"
	"strings"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/unstableunicorn/cogo/lib"
)

var clientMetadata []string
var userAttributes []string
var userGroups []string

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

func addUserToGroup(userName string, groupName string) {
	addUserToGroupInput := &cognito.AdminAddUserToGroupInput{
		UserPoolId: &poolID,
		Username:   &userName,
		GroupName:  &groupName,
	}

	_, err := cognitoSvc.AdminAddUserToGroup(addUserToGroupInput)

	if err != nil {
		lib.HandleAWSError("adding user to group", err, true)
	}
}

func getGroupsInUser(userName string) cognito.AdminListGroupsForUserOutput {
	listGroupsInUserInput := &cognito.AdminListGroupsForUserInput{
		UserPoolId: &poolID,
		Username:   &userName,
		Limit:      &limit,
	}

	var groups cognito.AdminListGroupsForUserOutput
	for {
		g, err := cognitoSvc.AdminListGroupsForUser(listGroupsInUserInput)

		if len(g.Groups) > 0 {
			for _, v := range g.Groups {
				groups.SetGroups(append(groups.Groups, v))
			}
		} else {
			groups = *g
		}

		if err != nil {
			lib.HandleAWSError("listing groups in user", err, true)
		}

		if g.NextToken == nil {
			break
		}

		listGroupsInUserInput.SetNextToken(*g.NextToken)
	}

	return groups
}
