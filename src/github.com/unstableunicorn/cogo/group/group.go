package group

import (
	"flag"
	"fmt"
	"os"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/unstableunicorn/cogo/app"
)

// Command parameters for group calls.
type Command struct {
	appOptions *app.App
	limit      int64
	nextToken  string
}

// Run command parameters.
func Run(opt *app.App) {
	groupOpts := flag.NewFlagSet("group options", flag.ExitOnError)
	limit := groupOpts.Int64("limit", 60, "The maximum amount of records to return")
	list := flag.NewFlagSet("list", flag.ExitOnError)

	if len(opt.Args) < 2 {
		fmt.Println("Expect some args: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	groupOpts.Parse(opt.Args[1:])
	command := opt.Args[len(opt.Args)-groupOpts.NArg():]
	switch command[0] {
	case "list", "ls":
		list.Parse(command[1:])
		cmd := &Command{
			appOptions: opt,
			limit:      *limit,
		}
		List(cmd)
	default:
		fmt.Println("Arguments not provide, please see the below")
		groupOpts.PrintDefaults()
		os.Exit(1)
	}
}

// List Cognito Groups in a given userpool.
func List(cmd *Command) {

	listGroupInput := &cognito.ListGroupsInput{
		Limit:      &cmd.limit,
		UserPoolId: &cmd.appOptions.UserPoolID,
	}

	users, err := cmd.appOptions.CognitoClient.ListGroups(listGroupInput)

	if err != nil {
		fmt.Println("Error getting groups", err)
	}

	fmt.Println(users.GoString())
}
