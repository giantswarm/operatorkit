package transaction

import (
	"strings"
)

func responseKey(keys ...string) string {
	return strings.Join(append([]string{"transaction", "responder"}, keys...), "/")
}

func transactionKey(keys ...string) string {
	return strings.Join(append([]string{"transaction", "executer"}, keys...), "/")
}
