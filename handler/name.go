package handler

import (
	"fmt"
	"strings"
)

func Name(r Interface) string {
	split := strings.Split(fmt.Sprintf("%#v", r), ".")

	if len(split) < 2 {
		return "n/a"
	}

	return strings.Replace(split[0], "&", "", 1)
}
