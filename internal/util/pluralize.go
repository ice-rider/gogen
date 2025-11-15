package util

import (
	"strings"
	"unicode"
)

func Pluralize(word string) string {
	if word == "" {
		return word
	}

	irregulars := map[string]string{
		"person": "people",
		"child":  "children",
		"man":    "men",
		"woman":  "women",
		"tooth":  "teeth",
		"foot":   "feet",
		"mouse":  "mice",
		"goose":  "geese",
	}

	lower := strings.ToLower(word)
	if plural, ok := irregulars[lower]; ok {

		if unicode.IsUpper(rune(word[0])) {
			return strings.ToUpper(plural[:1]) + plural[1:]
		}
		return plural
	}

	if strings.HasSuffix(lower, "s") ||
		strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") {
		return word + "es"
	}

	if strings.HasSuffix(lower, "y") {

		if len(word) > 1 && !isVowel(rune(lower[len(lower)-2])) {
			return word[:len(word)-1] + "ies"
		}
		return word + "s"
	}

	if strings.HasSuffix(lower, "f") {
		return word[:len(word)-1] + "ves"
	}

	if strings.HasSuffix(lower, "fe") {
		return word[:len(word)-2] + "ves"
	}

	if strings.HasSuffix(lower, "o") {

		if len(word) > 1 && !isVowel(rune(lower[len(lower)-2])) {
			return word + "es"
		}
		return word + "s"
	}

	return word + "s"
}

func Singularize(word string) string {
	if word == "" {
		return word
	}

	irregulars := map[string]string{
		"people":   "person",
		"children": "child",
		"men":      "man",
		"women":    "woman",
		"teeth":    "tooth",
		"feet":     "foot",
		"mice":     "mouse",
		"geese":    "goose",
	}

	lower := strings.ToLower(word)
	if singular, ok := irregulars[lower]; ok {
		if unicode.IsUpper(rune(word[0])) {
			return strings.ToUpper(singular[:1]) + singular[1:]
		}
		return singular
	}

	if strings.HasSuffix(lower, "ies") {
		return word[:len(word)-3] + "y"
	}

	if strings.HasSuffix(lower, "ves") {
		return word[:len(word)-3] + "f"
	}

	if strings.HasSuffix(lower, "ses") ||
		strings.HasSuffix(lower, "xes") ||
		strings.HasSuffix(lower, "zes") ||
		strings.HasSuffix(lower, "ches") ||
		strings.HasSuffix(lower, "shes") {
		return word[:len(word)-2]
	}

	if strings.HasSuffix(lower, "s") && len(word) > 1 {
		return word[:len(word)-1]
	}

	return word
}

func isVowel(r rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, r)
}
