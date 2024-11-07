package utils

import "strconv"

func IsSixDigitCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	_, err := strconv.Atoi(code)
	return err == nil
}
