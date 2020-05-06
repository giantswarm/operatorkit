package handler

import (
	"fmt"
	"strings"
)

func Name(r Interface) string {
	split := strings.Split(fmt.Sprintf("%#v", r), ".")

	if len(split) < 2 {
		panic("unable to parse handler name")
	}

	return strings.Replace(split[0], "&", "", 1)
}
