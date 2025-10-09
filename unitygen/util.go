package unitygen

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

func parseString(root *gabs.Container, path string) (string, error) {
	tagNameNode := root.Path(path)
	if tagNameNode == nil {
		return "", fmt.Errorf("%s is missing", path)
	}

	tagName, ok := tagNameNode.Data().(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string", path)
	}

	return tagName, nil
}

func pathToCsharpIdentifier(str string) string {
	result := strings.Builder{}
	capitalizeNext := true
	for _, char := range str {
		if char == '/' || char == '{' || char == '}' {
			capitalizeNext = true
			continue
		}

		if capitalizeNext {
			result.WriteString(strings.ToUpper(string(char)))
		} else {
			result.WriteRune(char)
		}
		capitalizeNext = false
	}
	return result.String()
}
