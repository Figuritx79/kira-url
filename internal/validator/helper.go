package validator

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

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
	return utf8.RuneCountInString(value) >= min
}

func MaxRunes(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}
