/*Package lib error handling
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
package lib

import (
	"fmt"
	"log"
	"strings"
)

// HandleAWSError strips and logs strings good format
func HandleAWSError(action string, err error, stopOnError bool) {
	_, extractedError := GetAWSError(err)
	errorString := fmt.Sprintf("Error %v,%v\n", action, extractedError)
	if stopOnError {
		log.Fatal(errorString)
	} else {
		log.Print(errorString)
	}
}

// GetAWSError strips and returns the error string
func GetAWSError(err error) (string, string) {
	errorStrings := strings.Split(err.Error(), ":")
	return errorStrings[0], errorStrings[1]
}
