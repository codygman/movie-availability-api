package api

import (
	"strings"
)

// This function takes two titles and returns the percent of match they are.
func MatchConfidence(title, entryTitle string) (float64) {
	// ??? Maybe ??? TODO: This could cause another array error in the case of "Entry" being empty, make sure it fails gracefully
	// TODO: Return error on failure
	matches :=0
	title = strings.ToLower(title)
	entryTitle = strings.ToLower(entryTitle)

	substrings := strings.Split(title, " ")
	// iterate over every substring and see if it is in the entry string, record matches
	for _, word := range substrings{
		if strings.Contains(strings.ToLower(word), title) {
			matches += 1
		}
	}

	percentOfMatches := float64((matches/len(substrings))*100.0)
	return percentOfMatches
}


