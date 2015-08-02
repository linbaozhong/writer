package utils

import (
	"fmt"
	"strings"
)

// 首填充
func StringStartPad(s, p string, l int) string {
	p = strings.Repeat(p, l-len(s))
	return fmt.Sprintf("%s%s", p, s)
}

// 尾填充
func StringEndPad(s, p string, l int) string {
	p = strings.Repeat(p, l-len(s))
	return fmt.Sprintf("%s%s", s, p)
}
