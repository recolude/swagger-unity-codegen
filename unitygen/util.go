package unitygen

import (
	"fmt"

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
