package txt

import (
	"fmt"
	"strings"
	"unicode"
)

func SelectLine(lines []string, lineIdx int) (string, error) {
	if len(lines) <= lineIdx {
		return "", fmt.Errorf("data needs at least %d line(s)", lineIdx+1)
	}
	return lines[lineIdx], nil
}

func stripFromFirstChar(s, chars string) string {
	if cut := strings.IndexAny(s, chars); cut >= 0 {
		return strings.TrimRightFunc(s[:cut], unicode.IsSpace)
	}
	return s
}
