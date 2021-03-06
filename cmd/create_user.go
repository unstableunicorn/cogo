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
	"errors"
	"fmt"

	"github.com/unstableunicorn/cogo/lib"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

var deliveryMedium []string
var tempPassword string

var suppress bool
var updateUserIfExists bool
var generatePasswordFlag bool

// createUserCmd represents the user command
var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create an AWS Cognito user",
	Long: `Create an AWS Cognito user.

  Examples: 
  Basic use, create a user named user1
  >cogo -p <poolid> create user user1

  Create a user with verified email and set password
  >cogo -p <poolid> create user user1 --password SuperSecretPassword --attributes "email=user1@company.com,email_verified=true"
  .`,
	Aliases: userAliases,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkUserArgs(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		userInput := createCreateUsersInput(args[0])
		createUser(userInput)
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringSliceVar(&clientMetadata, "clientmetadata", []string{}, "Comma separated list of client metadata e.g Key=Value[,Key2=Value2]")
	createUserCmd.Flags().StringSliceVar(&deliveryMedium, "delivery", []string{}, "List of desired delivery mediums, must be 'EMAIL' and/or 'SMS'")
	createUserCmd.Flags().BoolVar(&suppress, "suppress", false, "If set the user will NOT be send an email, default will send an email")
	createUserCmd.Flags().StringVar(&tempPassword, "password", "", "Password to set for the user, autogenerated if not provided, (Minimum Length: 6)")
	createUserCmd.Flags().StringSliceVar(&userAttributes, "attributes", []string{}, `Comma separated list of user attributes e.g Key=Value[,Key2=Value2]
 This is where the email is set. See usage`)
	createUserCmd.Flags().StringSliceVarP(&userGroups, "groups", "g", []string{}, "Add groups to the user as comma separated list")

	createUserCmd.Flags().BoolVarP(&updateUserIfExists, "update", "u", false, "If the user exists this will update the user")
	createUserCmd.Flags().BoolVar(&generatePasswordFlag, "autopassword", false, "Auto generate a password for the user")
}

func checkUserNameArg(args []string) error {
	if len(args) != 1 {
		s := fmt.Sprintf("An single value is required for the user name, provided values: '%v'", args)
		return errors.New(s)
	}
	return nil
}

func checkUserArgs(cmd *cobra.Command, args []string) error {
	hasUserError := checkUserNameArg(args)
	if hasUserError != nil {
		return hasUserError
	}

	if len(tempPassword) > 0 {
		if len(tempPassword) < 6 {
			return errors.New("Password does not meet minimum length")
		}
	} else {
		if generatePasswordFlag {
			pass, err := lib.GeneratePassword(12, 3, 3, 3)
			if err != nil {
				return err
			}
			tempPassword = pass
		} else {
			return errors.New("Password must be set or use --autopassword")
		}
	}

	if len(deliveryMedium) > 0 {
		for _, v := range deliveryMedium {
			if v != cognito.DeliveryMediumTypeEmail && v != cognito.DeliveryMediumTypeSms {
				return errors.New("Delivery medium must be EMAIL and or SMS")
			}
		}
	}

	return nil
}

func createCreateUsersInput(userName string) *cognito.AdminCreateUserInput {
	createUsersInput := &cognito.AdminCreateUserInput{
		UserPoolId: &poolID,
		Username:   &userName,
	}

	// Generate the mapped attributes
	mappedClientMetadata := mapClientMetadata(clientMetadata)

	if len(mappedClientMetadata) > 0 {
		createUsersInput.SetClientMetadata(mappedClientMetadata)
	}

	if len(deliveryMedium) > 0 {
		var dm []*string
		for _, v := range deliveryMedium {
			dm = append(dm, &v)
		}
		createUsersInput.SetDesiredDeliveryMediums(dm)
	}

	if updateUserIfExists {
		action := cognito.MessageActionTypeResend
		if suppress {
			action = cognito.MessageActionTypeSuppress
		}
		createUsersInput.SetMessageAction(action)
	}

	mappedUserAttributes := mapUserAttributes(userAttributes)
	if len(mappedUserAttributes) > 0 {
		createUsersInput.SetUserAttributes(mappedUserAttributes)
	}

	return createUsersInput
}

func createUser(userInput *cognito.AdminCreateUserInput) {
	user, err := cognitoSvc.AdminCreateUser(userInput)

	if err != nil {
		// try and update the user
		if updateUserIfExists {
			updateUser(&cognito.AdminUpdateUserAttributesInput{
				UserPoolId:     userInput.UserPoolId,
				Username:       userInput.Username,
				ClientMetadata: userInput.ClientMetadata,
				UserAttributes: userInput.UserAttributes,
			})
		} else {
			lib.HandleAWSError("creating user", err, true)
		}
	}

	if len(userGroups) > 0 {
		for _, group := range userGroups {
			addUserToGroup(*userInput.Username, group)
		}
	}

	fmt.Println(user.GoString())
}
