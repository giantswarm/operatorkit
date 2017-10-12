package metricsresource

import (
	"bytes"
	"regexp"
)

var camelCaseRegex = regexp.MustCompile("[0-9A-Za-z]+")

func toCamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelCaseRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
