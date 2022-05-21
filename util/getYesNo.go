package util

import (
	"fmt"
	"strings"
)

func GetYesNo() bool {
	for true {
		var userInput string
		_, err := fmt.Scanln(&userInput)
		if err != nil {
			if err.Error() != "unexpected newline" {
				panic(err)
			}
			continue
		}

		trimmed := strings.Trim(userInput, " ")
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
