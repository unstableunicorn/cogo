/*Package lib pass stdin arguments
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
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// GetArgsFromStdin allows arguments to be piped in
func GetArgsFromStdin() ([]string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("Not a pipe")
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	// Get the input from the pipe
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	if len(output) > 0 {
		args := stripSplitRune(output)
		if args != nil {
			return stripSplitRune(output), nil
		}
	}

	return nil, errors.New("No data from pipe")
}

func stripSplitRune(r []rune) []string {
	inString := string(r)
	noNewLines := strings.Replace(inString, "\r\n", " ", -1)
	noNewLines = strings.Replace(noNewLines, "\n", " ", -1)
	splitString := strings.Split(noNewLines, " ")

	//Remove spaces from each string
	var returnStrings []string
	for _, v := range splitString {
		var s string
		s = strings.TrimSpace(v)
		if len(s) > 0 {
			returnStrings = append(returnStrings, s)
		}
	}
	return returnStrings
}
