package convention

import "unicode"

func TitleCase(in string) string {
	if in == "" {
		return ""
	}
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
