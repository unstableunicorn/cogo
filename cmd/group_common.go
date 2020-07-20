/*Package cmd functions for group commands.
Copyright © 2020 Elric Hindy <anunstableunicorn@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
	"github.com/unstableunicorn/cogo/lib"
)

var precedence int64
var roleArn string
var groupDescription string

func checkGroupArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		s := fmt.Sprintf("An single value is required for the group name, provided values: '%v'", args)
		return errors.New(s)
	}

	if precedence < 0 {
		return errors.New("Precedence must be a positive value")
	}

	if len(roleArn) > 0 && len(roleArn) <= 20 {
		return errors.New("RoleArn must be greater than 20 characters")
	}
	return nil
}

func getUsersInGroup(groupName string) cognito.ListUsersInGroupOutput {
	listUsersInGroupInput := &cognito.ListUsersInGroupInput{
		UserPoolId: &poolID,
		GroupName:  &groupName,
		Limit:      &limit,
	}

	var users cognito.ListUsersInGroupOutput
	for {
		u, err := cognitoSvc.ListUsersInGroup(listUsersInGroupInput)

		if len(u.Users) > 0 {
			for _, v := range u.Users {
				users.SetUsers(append(users.Users, v))
			}
		} else {
			users = *u
		}

		if err != nil {
			lib.HandleAWSError("listing users in group", err, true)
		}

		if u.NextToken == nil {
			break
		}

		listUsersInGroupInput.SetNextToken(*u.NextToken)
	}

	return users
}
