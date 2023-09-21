package parse

import "strings"

/*
util.go

Contains any miscellaneous methods that aren't really related to
job parsing
*/

/*
hasKeyword

will split a string by it's spaces and see if it contains any keyword
*/
func hasKeyword(s string, keywords []string) bool {
	// iterating through each word in the title
	for _, token := range strings.Split(s, " ") {
		for _, keyword := range keywords {
			if strings.EqualFold(token, keyword) {
				return true
			}
		}
	}
	return false
}
