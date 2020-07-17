/*Package cmd root functions.
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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cfgFile string
var poolID string
var cognitoSvc *cognito.CognitoIdentityProvider

var groupAliases = []string{"g", "grp", "groups"}
var userAliases = []string{"u", "usr", "users"}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cogo",
	Short: "Cogo is a small cli utility to manage cognito users and groups",
	Long: `Usage: cogo [OPTIONS] [COMMAND]
  Cogo (short for Cognito Go)  is a cli written in Go that allows
  you to create, update, list and delete cognito users and groups including
  filtering and providing the ability to bulk update users.
  
  Example Usages:
  To list users:
  >cogo list users
  
  To list groups:
  >cogo list groups
  
  To list users and only show certain attributes:
  >cogo list users --attr username email status custom:somecustomattribute
  
  To create a user with sane defaults and add to existing groups:
  >cogo add user first.last@organisation.com --groups grp1 grp2

  Shortcuts! to make life easier you can use the following aliases:
  list|ls
  users|user|usr|u
  groups|group|grp|g
  
  e.g. cogo list users -> cogo ls u
  `,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !printVersionFlag && len(poolID) == 0 {
			return errors.New(`Required flag(s) "poolid" not set`)
		}

		if printVersionFlag {
			printVersion()
			os.Exit(0)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defer initCognito()
	rootCmd.PersistentFlags().StringVarP(&poolID, "poolid", "p", "", "AWS Cognito User PoolID (required)")
	rootCmd.PersistentFlags().BoolVarP(&printVersionFlag, "version", "v", false, "Print the current version string")
}

func initCognito() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		fmt.Println("Error creating session ", err)
	}

	cognitoSvc = cognito.New(sess)
}
