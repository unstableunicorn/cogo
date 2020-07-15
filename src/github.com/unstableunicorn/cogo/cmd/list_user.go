/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var listUserCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"u", "usr", "users"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		listUsers()
	},
}

func init() {
	listCmd.AddCommand(listUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listUsers() {
	listUsersInput := &cognito.ListUsersInput{
		Limit:      &limit,
		UserPoolId: &poolID,
	}

	for {
		users, err := cognitoSvc.ListUsers(listUsersInput)

		if err != nil {
			fmt.Println("Error getting users", err)
		}

		fmt.Println(users.GoString())

		if users.PaginationToken == nil {
			break
		}

		listUsersInput.PaginationToken = users.PaginationToken
	}
}
