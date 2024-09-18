package firebase

import (
	"fmt"
	"regexp"
	"strconv"
)

// describe encodeEmail Each illegal character `[.#$\[\]]` is replaced with an underscore followed by its two-digit hexadecimal ASCII code.
func encodeEmail(email string) string {
	pattern := regexp.MustCompile(`[.#$\[\]]`)

	return pattern.ReplaceAllStringFunc(email, func(m string) string {
		return "_" + fmt.Sprintf("%X", m[0])
	})
}

/*
	describe decodeEmail The function finds all sequences that match the pattern _(\w{2}).

It then converts the hexadecimal code back to the original character.
*/
func decodeEmail(encodedEmail string) string {
	pattern := regexp.MustCompile(`_(\w{2})`)

	return pattern.ReplaceAllStringFunc(encodedEmail, func(m string) string {
		hexCode := m[1:] // Remove the underscore
		charCode, err := strconv.ParseInt(hexCode, 16, 32)
		if err != nil {
			// Handle error (you may want to log or return an error)
			return m
		}
		return string(rune(charCode))
	})
}
