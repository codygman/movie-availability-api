package api

import (
	"strings"
)


func normalizeKeyword(keyword string) (string) {
	// normalizes keyword for comparison
	// do your own escaping!
	var replacementMap = make(map[string]string)
	replacementMap["·"] = "-"

	for old, replacement := range replacementMap {
		keyword = strings.Replace(keyword, old, replacement, -1)
	} 
	return keyword

}
