package convention

import "unicode"

func TitleCase(in string) string {
	if in == "" {
		return ""
	}

	out := make([]rune, 0)

	nextCapitilized := true
	for _, c := range in {
		if c == '_' || c == '-' {
			nextCapitilized = true
			continue
		}

		if nextCapitilized {
			out = append(out, unicode.ToUpper(c))
			nextCapitilized = false
		} else {
			out = append(out, c)
		}
	}

	return string(out)
}
