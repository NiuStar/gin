package utils

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Convert2String(value interface{}) string {
	switch value.(type) {
	case int:
		{
			return strconv.FormatInt(int64(value.(int)), 10)
		}
		break
	case int64:
		{
			return strconv.FormatInt(value.(int64), 10)
		}
		break
	case float32:
		{
			return strconv.FormatFloat(float64(value.(float32)), 'f', 0, 32)
		}
		break
	case float64:
		{
			return strconv.FormatFloat(value.(float64), 'f', 0, 64)
		}
		break
	case string:
		{
			return value.(string)
		}
		break
	case []byte:
		{
			return string(value.([]byte))
		}
		break
	}
	return ""
}

func ToUpperOnlyFirstChar(str string) string {

	s := str[:1]

	isASCII, hasLower := true, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasLower = hasLower || (c >= 'a' && c <= 'z')
	}

	if isASCII { // optimize for ASCII-only strings.
		if !hasLower {
			return s
		}
		b := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
			b[i] = c
		}
		return string(b) + str[1:]
	}
	return strings.Map(unicode.ToUpper, s) + str[1:]
}

func SnackCase(raw string) string {
	newstr := make([]rune, 0, len(raw))
	for i, chr := range raw {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if i > 0 {
				newstr = append(newstr, '_')
			}
		}
		newstr = append(newstr, unicode.ToLower(chr))
	}
	return string(newstr)
}
