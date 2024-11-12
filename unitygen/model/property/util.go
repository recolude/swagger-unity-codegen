package property

import (
	"fmt"
	"strings"
)

func WriteDescriptionComment(builder *strings.Builder, description string) {
	cleaned := strings.TrimSpace(description)
	if cleaned != "" {
		builder.WriteString("\t/// <summary>\n")
		fmt.Fprintf(builder, "\t/// %s\n", cleaned)
		builder.WriteString("\t/// </summary>\n")
	}
}
