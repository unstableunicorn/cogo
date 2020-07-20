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
