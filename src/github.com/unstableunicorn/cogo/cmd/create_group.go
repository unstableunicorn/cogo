/*Package cmd functions for create group commands.
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

var updateGroupIfExists bool

// createGroupCmd represents the group command
var createGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create an AWS Cognito group",
	Long: `Create an AWS Cognito group.

  Examples: 
  Basic use, create a group named 'mynewgroup'
  >cogo -p <poolid> create group mynewgroupname

  Create a group with description and precedence
  >cogo -p <poolid> create group Group.Name -d "A group that does group things" --precedence 3 
  .`,
	Aliases: groupAliases,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkArgs(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		groupInput := createCreateGroupsInput(args[0])
		createGroup(groupInput)
	},
}

func init() {
	createCmd.AddCommand(createGroupCmd)

	createGroupCmd.Flags().Int64Var(&precedence, "precedence", 0, "Group precedence")
	createGroupCmd.Flags().StringVarP(&groupDescription, "description", "d", "", "Description of the group")
	createGroupCmd.Flags().StringVarP(&roleArn, "role", "r", "", "Description of the group")
	createGroupCmd.Flags().BoolVarP(&updateGroupIfExists, "update", "u", false, "If the group exists this will update the group")
}

func createCreateGroupsInput(groupName string) *cognito.CreateGroupInput {
	createGroupsInput := &cognito.CreateGroupInput{
		UserPoolId: &poolID,
		GroupName:  &groupName,
	}

	if len(groupDescription) > 0 {
		createGroupsInput.Description = &groupDescription
	}

	if len(roleArn) > 0 {
		createGroupsInput.RoleArn = &roleArn
	}

	if precedence != 0 {
		createGroupsInput.Precedence = &precedence
	}

	return createGroupsInput
}

func createGroup(groupInput *cognito.CreateGroupInput) {
	group, err := cognitoSvc.CreateGroup(groupInput)

	if err != nil {
		// try and update the group
		if updateGroupIfExists {
			ug := cognito.UpdateGroupInput(*groupInput)
			updateGroup(&ug)
		} else {
			lib.HandleAWSError("creating group", err)
		}
	} else {
		fmt.Println(group.GoString())
	}
}
