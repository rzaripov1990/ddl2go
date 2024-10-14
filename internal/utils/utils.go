package utils

import (
	"fmt"
	"strings"
)

func ToCamelCase(input string) (upperCS, classicCS string) {
	var (
		foundUnderscore bool
		countUnderscore int
	)

	for i := range input {
		if input[i] == '_' {
			countUnderscore++
			if !foundUnderscore {
				foundUnderscore = true
			}
		} else {
			break
		}
	}
	if foundUnderscore {
		input = strings.Replace(input, strings.Repeat("_", countUnderscore), fmt.Sprintf("Underscore_%d", countUnderscore), 1)
	}

	var words []string = strings.Split(input, "_")
	for i := range words {
		if words[i] != "" {
			lower := strings.ToLower(words[i][1:])
			upper := strings.ToUpper(words[i][:1])

			words[i] = upper + lower
		}
	}
	upperCS = strings.Join(words, "")
	classicCS = strings.ToLower(upperCS[:1]) + upperCS[1:]
	return
}
