package util

import (
	"errors"
	"regexp"
)

var (
	// ErrNoMatchFound ...
	ErrNoMatchFound = errors.New("No match found")
)

// RegexExtractMatches extracts multiple named matches from a string using regex
func RegexExtractMatches(in, regExp string, names ...string) (map[string]string, error) {
	compiledRegex, err := regexp.Compile(regExp)
	if err != nil {
		return map[string]string{}, err
	}

	subExpNames := compiledRegex.SubexpNames()
	results := compiledRegex.FindAllStringSubmatch(in, -1)

	if len(results) == 0 {
		return map[string]string{}, ErrNoMatchFound
	}

	matches := make(map[string]string, len(names))
	for i, match := range results[0] {
		for _, name := range names {
			if subExpNames[i] == name {
				matches[name] = match
			}
		}
	}

	return matches, nil
}

// RegexExtractMatch extracts a named match from a string using regex
func RegexExtractMatch(in, regularExpression, name string) (string, error) {
	compiledRegex, err := regexp.Compile(regularExpression)
	if err != nil {
		return "", err
	}

	names := compiledRegex.SubexpNames()
	results := compiledRegex.FindAllStringSubmatch(in, -1)

	if len(results) == 0 {
		return "", ErrNoMatchFound
	}

	for i, match := range results[0] {
		if names[i] != name {
			continue
		}

		return match, nil
	}

	return "", ErrNoMatchFound
}
