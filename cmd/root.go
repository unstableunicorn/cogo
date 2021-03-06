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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unstableunicorn/cogo/lib"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cfgFile string
var poolID string
var cognitoSvc *cognito.CognitoIdentityProvider
var version = "developement.version"

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
  
  Examples:
  To list users:
  >cogo -p <poolid> list users
  
  To list groups:
  >cogo -p <poolid> list groups
  
  To list users and only show certain attributes:
  >cogo -p <poolid> list users --attr username email status custom:somecustomattribute
  
  To create a user with sane defaults and add to existing groups:
  >cogo -p <poolid> add user first.last@organisation.com --groups grp1 grp2

  Shortcuts! to make life easier you can use the following aliases:
  list|ls
  users|user|usr|u
  groups|group|grp|g
  e.g. cogo -p <poolid> list users -> cogo ls u

  You can also enter the poolid anywhere on the command yay:
  >cogo create user username -email -p <poolid>
  >cogo create user -p <poolid> username -email   
  >cogo create -p <poolid> user username -email   
  `,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Get args from pipe and append to existing args
	pipedArgs, _ := lib.GetArgsFromStdin()
	if len(pipedArgs) > 0 {
		allArgs := os.Args[1:]
		for _, arg := range pipedArgs {
			allArgs = append(allArgs, arg)
		}
		rootCmd.SetArgs(allArgs)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defer initCognito()
	rootCmd.PersistentFlags().StringVarP(&poolID, "poolid", "p", "", "AWS Cognito User PoolID (required)")
	rootCmd.MarkPersistentFlagRequired("poolid")
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
