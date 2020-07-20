/*Package cmd functions for create user commands.
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
	"fmt"

	"github.com/unstableunicorn/cogo/lib"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

// updateUserCmd represents the user command
var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Update an AWS Cognito user",
	Long: `Update an AWS Cognito user.

  Examples: 
  Basic use, Update a user named user1
  >cogo -p <poolid> create user user1 --attributes "email=user1@newemail.com,email_verified=true"
  `,
	Aliases: userAliases,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkUserNameArg(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		userInput := createUpdateUsersInput(args[0])
		updateUser(userInput)
	},
}

func init() {
	updateCmd.AddCommand(updateUserCmd)

	updateUserCmd.Flags().StringSliceVar(&clientMetadata, "clientmetadata", []string{}, "Comma separated list of client metadata e.g Key=Value[,Key2=Value2]")
	updateUserCmd.Flags().StringSliceVar(&userAttributes, "attributes", []string{}, `Comma separated list of user attributes e.g Key=Value[,Key2=Value2]
 This is where the email is set. See usage`)
}

func createUpdateUsersInput(userName string) *cognito.AdminUpdateUserAttributesInput {
	updateUserInput := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: &poolID,
		Username:   &userName,
	}

	// Generate the mapped attributes
	mappedClientMetadata := mapClientMetadata(clientMetadata)

	if len(mappedClientMetadata) > 0 {
		updateUserInput.SetClientMetadata(mappedClientMetadata)
	}

	mappedUserAttributes := mapUserAttributes(userAttributes)
	if len(mappedUserAttributes) > 0 {
		updateUserInput.SetUserAttributes(mappedUserAttributes)
	}

	return updateUserInput
}

func updateUser(userInput *cognito.AdminUpdateUserAttributesInput) {
	_, err := cognitoSvc.AdminUpdateUserAttributes(userInput)

	if err != nil {
		lib.HandleAWSError("creating user", err, true)
	} else {
		fmt.Println("Successfully updated user", *userInput.Username)
	}
}
