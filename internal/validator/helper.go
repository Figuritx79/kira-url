package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var regexCustomCode = regexp.MustCompile("^[a-z0-9]+(?:[-_][a-z0-9]+)*$")

func IsURL(value string) bool {
	u, err := url.Parse(value)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func NotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinRunes(value string, min int) bool {
	fmt.Println("value", value, "min", min, "runes", utf8.RuneCountInString(value))
	return utf8.RuneCountInString(value) >= min
}

func MaxRunes(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func IsValidCustomCode(value string) bool {
	return regexCustomCode.MatchString(value)
}
