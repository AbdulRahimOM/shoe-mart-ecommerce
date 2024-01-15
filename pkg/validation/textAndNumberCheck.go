package requestValidation

import (
	"strconv"
	"unicode"
)

func IsValidTextOrNumber(s string, minLength int) bool {
	// Check if the string contains at least minLength alphanumeric characters
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
			if count >= minLength {
				return true
			}
		}
	}
	return false
}

func IsValidText(s string, minLength int) bool {
	// Check if the string contains at least minLength alphanumeric characters
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
			if count >= minLength {
				return true
			}
		}
	}
	return false
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func ContainsOnlyAlphabets(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func ContainsOnlyAlphabetsAndSpaces(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

func ContainsOnlyAlphabetsDigitsAndSpaces(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}