package utils

import (
	"strings"
	"unicode"
)

func ContainsUppercase(s string) bool {
	for i, n := 0, len(s); i < n; i++ {
		if unicode.IsUpper(rune(s[i])) {
			return true
		}
	}

	return false
}

func ContainsLowercase(s string) bool {
	for i, n := 0, len(s); i < n; i++ {
		if unicode.IsLower(rune(s[i])) {
			return true
		}
	}

	return false
}

func ContainsNumber(s string) bool {
	for i, n := 0, len(s); i < n; i++ {
		if unicode.IsNumber(rune(s[i])) {
			return true
		}
	}

	return false
}

func ContainsSpecialChar(s string) bool {
	for i, n := 0, len(s); i < n; i++ {
		if strings.Contains(SPECIAL_CHARS, string(s[i])) {
			return true
		}
	}

	return false
}

func ContainsWhitespace(s string) bool {
	for i, n := 0, len(s); i < n; i++ {
		if unicode.IsSpace(rune(s[i])) {
			return true
		}
	}

	return false
}
