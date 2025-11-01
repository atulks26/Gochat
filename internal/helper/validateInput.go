package helper

import "strings"

func ValidateInput(input string, cmd string, parts int) (bool, []string) {
	divided := strings.SplitN(strings.TrimSpace(input), " ", parts)
	if len(divided) != parts || strings.ToUpper(divided[0]) != cmd {
		return false, nil
	}

	return true, divided
}
