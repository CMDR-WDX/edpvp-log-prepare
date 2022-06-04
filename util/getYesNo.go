package util

import (
	"bufio"
	"strings"
)

func GetYesNo(reader *bufio.Reader, defaultValue bool) bool {
	for true {
		var userInput string
		reader.Discard(reader.Buffered())
		userInput, err := reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\r\n")
		if err != nil {
			if err.Error() != "unexpected newline" {
				panic(err)
			}
			continue
		}

		trimmed := strings.Trim(userInput, " ")
		if len(trimmed) == 0 {
			return defaultValue
		}

		if len(trimmed) != 1 {
			continue
		}

		if trimmed[0] == 'n' || trimmed[0] == 'N' {
			return false
		} else if trimmed[0] == 'y' || trimmed[0] == 'Y' {
			return true
		}
	}

	// Will never happen bc of the infinite loop
	return false

}
