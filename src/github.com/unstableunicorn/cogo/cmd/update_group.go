/*Package cmd functions for update group commands.
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

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
	"github.com/unstableunicorn/cogo/lib"
)

var createGroupIfNotExists bool

// groupCmd represents the group command
var updateGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Update an AWS cognito group",
	Long: `Update an AWS Cognito group

  Examples: 
  Basic use, update a group named 'mynewgroup'
  >cogo -p <poolid> update group mynewgroupname -d "My updated description"

  Update a group with description and precedence
  >cogo -p <poolid> update group Group.Name -d "An updated group" --precedence 3 
  `,
	Aliases: groupAliases,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkGroupArgs(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		groupInput := createUpdateGroupInput(args[0])
		updateGroup(groupInput)
	},
}

func init() {
	updateCmd.AddCommand(updateGroupCmd)
	updateGroupCmd.Flags().Int64Var(&precedence, "precedence", 0, "Group precedence")
	updateGroupCmd.Flags().StringVarP(&groupDescription, "description", "d", "", "Description of the group")
	updateGroupCmd.Flags().StringVarP(&roleArn, "role", "r", "", "Description of the group")
	updateGroupCmd.Flags().BoolVarP(&createGroupIfNotExists, "create", "c", false, "If the group does not exists this will create the group")
}

func createUpdateGroupInput(groupName string) *cognito.UpdateGroupInput {
	updateGroupsInput := &cognito.UpdateGroupInput{
		UserPoolId: &poolID,
		GroupName:  &groupName,
	}

	if len(groupDescription) > 0 {
		updateGroupsInput.Description = &groupDescription
	}

	if len(roleArn) > 0 {
		updateGroupsInput.RoleArn = &roleArn
	}

	if precedence != 0 {
		updateGroupsInput.Precedence = &precedence
	}

	return updateGroupsInput
}

func updateGroup(groupInput *cognito.UpdateGroupInput) {
	updatedGroup, err := cognitoSvc.UpdateGroup(groupInput)

	if err != nil {
		if createGroupIfNotExists {
			cg := cognito.CreateGroupInput(*groupInput)
			createGroup(&cg)
		} else {
			lib.HandleAWSError("updating group", err)
		}

	} else {
		fmt.Println(updatedGroup.GoString())
	}
}
