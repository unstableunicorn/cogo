package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unstableunicorn/cogo/group"
	"github.com/unstableunicorn/cogo/user"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/unstableunicorn/cogo/app"
)

// CognitoCommand is a definition for all commands.
type CognitoCommand func(*app.App)

func main() {
	// Subcommands
	poolID := flag.String("poolid", "", "(Required) The ID of the user pool")
	// region := flag.String("region", "", "Specify the aws region, will user from profile if not specified")

	if len(os.Args) < 2 {
		fmt.Println("Expect some args")
		os.Exit(1)
	}

	flag.Parse()
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Error creating session ", err)
	}
	svc := cognito.New(sess)

	appInput := &app.App{
		CognitoClient: svc,
		UserPoolID:    *poolID,
		// Region:        *region,
		Args: os.Args[len(os.Args)-flag.NArg():],
	}

	switch appInput.Args[0] {
	case "grp", "group", "groups":
		group.Run(appInput)
	case "u", "user", "users":
		user.Run(appInput)
	default:
		flag.PrintDefaults()
	}

}
