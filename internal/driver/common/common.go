package common

import "strings"

// sanitizeTag returns the input string replacing all '/' or ':' symbols for an to an '_'(undersocre)
func SanitizeTag(input string) string {

	chars := map[string]string{
		"/": "_",
		":": "_",
	}

	for originalChar, newChar := range chars {
		input = strings.ReplaceAll(input, originalChar, newChar)
	}

	return input
}
