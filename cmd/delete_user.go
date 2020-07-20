/*Package cmd functions for delete user commands.
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

	"github.com/unstableunicorn/cogo/lib"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

// deleteUserCmd represents the user command
var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete an AWS Cognito user",
	Long: `Delete an AWS Cognito user, takes the user name(s)
  as an input.

  Example:
  Delete one user:
  >cogo -p <poolid> delete user user1

  Delete multiple users:
  >cogo -p <poolid> delete user user1name user2name
  `,

	Aliases: userAliases,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			s := fmt.Sprintf("At least one value is required for the user name, provided values: '%v'", args)
			return errors.New(s)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		deleteUser(args)
	},
}

func init() {
	deleteCmd.AddCommand(deleteUserCmd)
}

func deleteUser(usernames []string) {
	deleteUserInput := &cognito.AdminDeleteUserInput{
		UserPoolId: &poolID,
	}

	for _, username := range usernames {
		deleteUserInput.SetUsername(username)

		_, err := cognitoSvc.AdminDeleteUser(deleteUserInput)
		if err != nil {
			s := fmt.Sprintf("deleting user %v", username)
			lib.HandleAWSError(s, err, stopDeleteOnError)
		} else {
			fmt.Println("Successfully deleted: ", username)
		}
	}
}
