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

var precedence int64
var roleArn string

// createGroupCmd represents the group command
var createGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: groupAliases,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createGroup(args[0])
	},
}

func init() {
	createCmd.AddCommand(createGroupCmd)

	// Here you will define your flags and configuration settings.
	createGroupCmd.Flags().Int64Var(&precedence, "precedence", 3, "Group precedence")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createGroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createGroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createGroup(groupName string) {
	createGroupsInput := &cognito.CreateGroupInput{
		UserPoolId: &poolID,
		GroupName:  &groupName,
	}

	groups, err := cognitoSvc.CreateGroup(createGroupsInput)

	if err != nil {
		fmt.Println("Error creating group", err)
	}

	fmt.Println(groups.GoString())
}
