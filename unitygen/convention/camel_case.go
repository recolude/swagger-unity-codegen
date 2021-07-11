package convention

import "unicode"

// CamelCase translates string to start with lowercase
func CamelCase(in string) string {
	if in == "" {
		return ""
	}

	out := make([]rune, 0)

	nextCapitilized := false
	experiencedFist := false
	for _, c := range in {
		if experiencedFist == false {
			if c != '_' && c != '-' {
				out = append(out, unicode.ToLower(c))
				experiencedFist = true
			}
			continue
		}

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
