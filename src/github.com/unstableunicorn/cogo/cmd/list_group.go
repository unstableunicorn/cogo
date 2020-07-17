/*Package cmd functions for list group commands.
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

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
	"github.com/unstableunicorn/cogo/lib"
)

var filterGroupList string

// listGroupCmd represents the group command
var listGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Returns a list of groups",
	Long: `Returns a list of groups and provides additional filtering
  
  Filtering uses valid regex based on RE2 syntax, example:
  Return all groups that have admin in the name:
  >cogo -p <poolid> list groups -a -f "admin"
  
  Return all groups that start with admin or Admin:
  >cogo -p <poolid> list groups -a -f "^[aA]dmin"
  `,
	Aliases: groupAliases,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(filterGroupList) > 0 {
			_, err := regexp.Compile(filterGroupList)
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		listGroups()
	},
}

func init() {
	listCmd.AddCommand(listGroupCmd)
	listGroupCmd.Flags().StringVarP(&filterGroupList, "filter", "f", "", "A regex compatible string for listing groups")
}

func listGroups() {
	var matchGroupName *regexp.Regexp

	if len(filterGroupList) > 0 {
		matchGroupName = regexp.MustCompile(filterGroupList)
	}

	listGroupsInput := &cognito.ListGroupsInput{
		Limit:      &limit,
		UserPoolId: &poolID,
	}

	var groups cognito.ListGroupsOutput
	for {
		g, err := cognitoSvc.ListGroups(listGroupsInput)

		if len(g.Groups) > 0 {
			for _, v := range g.Groups {
				// Run the filter here
				if len(filterGroupList) > 0 {
					// If no match go to the next one
					if !matchGroupName.Match([]byte(*v.GroupName)) {
						continue
					}
				}
				groups.Groups = append(groups.Groups, v)
			}
		} else {
			groups = *g
		}

		if err != nil {
			lib.HandleAWSError("listing groups", err)
		}

		if groups.NextToken == nil || !getall {
			break
		}

		listGroupsInput.NextToken = g.NextToken
	}

	if len(groups.Groups) > 0 {
		fmt.Println(groups.GoString())
	} else {
		if len(filterGroupList) > 0 {
			errorString := fmt.Sprintf("No groups found matching filter: '%v'\n", filterGroupList)
			log.Fatal(errorString)
		} else {
			fmt.Println("No groups found")
		}
	}
}
