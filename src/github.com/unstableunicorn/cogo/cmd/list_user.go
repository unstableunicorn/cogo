/*Package cmd functions for list user commands.
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
	"log"
	"regexp"
	"strings"

	"github.com/unstableunicorn/cogo/lib"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

var listUserAttributes []string
var filterUserAttributes []string
var filterUserStatus string
var filterEnabledUsers bool

// userCmd represents the user command
var listUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Returns a list of users",
	Long: `Returns a list of users and provides additional filtering

  Filtering uses valid regex based on RE2 syntax, example:
  Return all users that have status CHANGE in the status field:
  cogo -p <poolid> list users -a --status "CHANGE"

  Return enabled users where their email has '@companyname':
  >cogo -p <poolid> list users --fattr "@companyname" 
     `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(filterUserStatus) > 0 {
			_, err := regexp.Compile(filterUserStatus)
			return err
		}
		return nil
	},
	Aliases: userAliases,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		listUsers()
	},
}

func init() {
	listCmd.AddCommand(listUserCmd)

	listUserCmd.Flags().StringSliceVar(&listUserAttributes, "attr", []string{}, "List of user attributes to return")
	listUserCmd.Flags().StringVar(&filterUserStatus, "status", "", "Return users with status that matches filter provided")
	listUserCmd.Flags().BoolVar(&filterEnabledUsers, "enabled", true, `Return users who are enabled, ignored if --all flag is set`)
	listUserCmd.Flags().StringSliceVar(&filterUserAttributes, "fattr", []string{}, "Key Value attribute pairs to filter in Key=Value[,Key2=Value2]")
}

func listUsers() {
	var matchUserStatus *regexp.Regexp

	if len(filterUserStatus) > 0 {
		matchUserStatus = regexp.MustCompile(filterUserStatus)
	}

	// Generate the mapped attributes
	mappedAttributeFilter := make(map[string]*regexp.Regexp)
	if len(filterUserAttributes) > 0 {
		for _, a := range filterUserAttributes {
			values := strings.Split(a, "=")
			if len(values) != 2 {
				errorString := fmt.Sprintf("Invalid attribute filter %v \n", a)
				log.Fatal(errorString)
			}
			mappedAttributeFilter[values[0]] = regexp.MustCompile(values[1])
		}
	}

	if getall && !filterEnabledUsers {
		fmt.Println("Ignoring enabled users filter")
	}

	listUsersInput := &cognito.ListUsersInput{
		Limit:      &limit,
		UserPoolId: &poolID,
	}

	// If attribute list is provided we need to convert the *[]string to a []*string
	if len(listUserAttributes) > 0 {
		for _, v := range listUserAttributes {
			listUsersInput.AttributesToGet = append(listUsersInput.AttributesToGet, &v)
		}
	}

	var users cognito.ListUsersOutput
	for {
		u, err := cognitoSvc.ListUsers(listUsersInput)

		if len(u.Users) > 0 {
			for _, v := range u.Users {
				// Run the filters here
				if len(filterUserStatus) > 0 {
					if !matchUserStatus.Match([]byte(*v.UserStatus)) {
						continue
					}
				}

				// Match Attribute values if defined
				if len(mappedAttributeFilter) > 0 {
					matchFound := false
					for _, attr := range v.Attributes {
						for k, v := range mappedAttributeFilter {
							if k == *attr.Name && v.Match([]byte(*attr.Value)) {
								matchFound = true
							}
						}
					}
					if !matchFound {
						continue
					}
				}

				if *v.Enabled != filterEnabledUsers && !getall {
					continue
				}
				users.Users = append(users.Users, v)
			}
		} else {
			users = *u
		}

		if err != nil {
			lib.HandleAWSError("listing users", err)
		}

		if users.PaginationToken == nil || !getall {
			break
		}

		listUsersInput.PaginationToken = users.PaginationToken
	}

	if len(users.Users) > 0 {
		fmt.Println(users.GoString())
	} else {
		fmt.Println("No users found matching provided arguments")
	}
}
