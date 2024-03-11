package main

import (
	"strings"
	"unicode/utf8"
)

// Pluralize function to convert a singular word to its plural form
func Pluralize(word string) string {
	if word == "" {
		return ""
	}

	// Simple rules for pluralization
	switch {
	case strings.HasSuffix(word, "y"):
		// city -> cities, baby -> babies
		return strings.TrimSuffix(word, "y") + "ies"
	case strings.HasSuffix(word, "s"),
		strings.HasSuffix(word, "sh"),
		strings.HasSuffix(word, "ch"),
		strings.HasSuffix(word, "x"),
		strings.HasSuffix(word, "z"):
		// bus -> buses, bush -> bushes, match -> matches
		return word + "es"
	default:
		// cat -> cats, dog -> dogs
		return word + "s"
	}
}

// CapitalizeFirst 将字符串的首个字母变成大写
func CapitalizeFirst(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(r)) + s[size:]
}

// ToCamelCase converts snake_case strings to CamelCase
func ToCamelCase(str string) string {
	var camelCase string
	words := strings.Split(str, "_")

	for _, word := range words {
		camelCase += strings.Title(word)
	}

	return camelCase
}

// ToLowerCamelCase converts snake_case strings to lowerCamelCase
func ToLowerCamelCase(str string) string {
	// Split the string by underscores
	words := strings.Split(str, "_")
	for i, word := range words {
		// Capitalize the first letter of each word except the first one
		if i != 0 {
			words[i] = strings.Title(word)
		}
	}
	// Join the words back together
	return strings.Join(words, "")
}
