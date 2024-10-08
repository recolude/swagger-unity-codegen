package convention

import (
	"bytes"
	"strings"
	"unicode"
)

func AddSpace(s string) string {
	buf := &bytes.Buffer{}

	capitalizeNext := false

	for i, rune := range strings.TrimSpace(s) {

		if i == 0 {
			buf.WriteRune(unicode.ToUpper(rune))
			continue
		}

		if rune == '-' || rune == '_' || rune == ' ' {
			buf.WriteRune(' ')
			capitalizeNext = true
			continue
		}

		if unicode.IsUpper(rune) {
			capitalizeNext = false
			buf.WriteRune(' ')
		}

		if capitalizeNext {
			buf.WriteRune(unicode.ToUpper(rune))
		} else {
			buf.WriteRune(rune)
		}

		capitalizeNext = false
	}
	return buf.String()
}
