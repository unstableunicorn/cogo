package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var precedence int64
var roleArn string
var groupDescription string

func checkArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		s := fmt.Sprintf("An single value is required for the group name, provided values: '%v'", args)
		return errors.New(s)
	}

	if precedence < 0 {
		return errors.New("Precedence must be a positive value")
	}

	if len(roleArn) > 0 && len(roleArn) <= 20 {
		return errors.New("RoleArn must be greater than 20 characters")
	}
	return nil
}
