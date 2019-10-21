package cmn

import (
	"bytes"
	"strings"
)

func CutTitle(str string, cutLen int) string {
	if len(str) <= cutLen {
		return str
	}
	return str[:cutLen] + "..."
}

func TrimLines(str string, cutLen int) string {
	var buf bytes.Buffer
	for _, v := range strings.Split(str, "\n") {
		buf.WriteString(CutTitle(v, cutLen) + "\n")
	}
	return buf.String()
}
