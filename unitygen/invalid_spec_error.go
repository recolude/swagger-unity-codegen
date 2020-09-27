package unitygen

import (
	"fmt"
	"strings"
)

// InvalidSpecError is an error around the structure of the specified swagger
// spec, hindering codegen capabilities
type InvalidSpecError struct {
	Path   []string
	Reason string
}

func (ise InvalidSpecError) Error() string {
	if len(ise.Path) == 0 {
		return fmt.Sprintf("Invalid spec: %s", ise.Reason)
	}

	var pathBuilder strings.Builder
	for i, str := range ise.Path {
		pathBuilder.WriteString(str)
		if i < len(ise.Path)-1 {
			pathBuilder.WriteRune('.')
		}
	}
	return fmt.Sprintf("Invalid spec at %s: %s", pathBuilder.String(), ise.Reason)
}
