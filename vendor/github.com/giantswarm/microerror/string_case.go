package microerror

import (
	"strings"
	"unicode"
)

func toStringCase(input string) string {
	chunks := []string{}
	split := strings.Split(input, "")

	for i, s := range split {
		r := []rune(s)

		var nextUpper bool
		if i != 0 && i+1 < len(split) {
			p := []rune(split[i-1])
			n := []rune(split[i+1])
			nextUpper = unicode.IsUpper(p[0]) && unicode.IsUpper(n[0])
		}

		isFirst := i == 0
		isLast := i+1 == len(split)
		isUpper := unicode.IsUpper(r[0])
		isAbbreviation := isUpper && nextUpper

		if !isAbbreviation && !isFirst && !isLast && isUpper {
			chunks = append(chunks, string(" "))
		}

		chunks = append(chunks, strings.ToLower(s))
	}

	return strings.Join(chunks, "")
}
