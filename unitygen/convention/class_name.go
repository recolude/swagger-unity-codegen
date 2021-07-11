package convention

import "unicode"

// ClassName cleans up a string to be a valid C# class name, removing invalid
// characters and capitilizing the first letter. Some of this is just stylistic
// choice.
func ClassName(in string) string {
	if in == "" {
		return ""
	}

	out := make([]rune, 0)

	nextCapitilized := false
	experiencedFist := false
	for _, c := range in {
		if experiencedFist == false {
			if c != '_' && c != '-' {
				out = append(out, unicode.ToUpper(c))
				experiencedFist = true
			}
			continue
		}

		if c == '-' {
			nextCapitilized = true
			continue
		}

		if c == '_' {
			out = append(out, c)
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
