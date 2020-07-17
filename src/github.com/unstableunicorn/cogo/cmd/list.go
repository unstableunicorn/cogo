/*
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

	"github.com/spf13/cobra"
)

var limit int64
var getall bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list cognito users or groups",
	Long: `Allows the user to list cognito users and groups
and provide filters for searching.

Usage: cogo list [group|user] [OPTIONS]
List the first 60 users with all attributes:
>cogo list users

List all groups:
>cogo list groups --all

List all users and only fetch the username and email fields:
>cogo list users --all --attribute username --attribute email

You can also combine attributes using a comma separated list:
>cogo list users --all --attr username,email,phone
`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().Int64VarP(&limit, "limit", "l", 60, "Maximum number of records to return [Max: 60]")
	listCmd.PersistentFlags().BoolVarP(&getall, "all", "a", false, "Return all the records, this will ignore the limit if set and other fields where listed in help")
}
