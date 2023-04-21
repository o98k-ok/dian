package builder

import (
	"fmt"
)

func Openfile(filename string, line int) string {
	return fmt.Sprintf("obsidian://advanced-uri?vault=obsidian&filepath=%s&line=%d", filename, line)
}
