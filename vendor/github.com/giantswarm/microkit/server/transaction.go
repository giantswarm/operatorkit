package server

import (
	"regexp"
)

const (
	// TransactionIDHeader is the canonical representation of the transaction ID
	// HTTP header field.
	TransactionIDHeader = "X-Transaction-ID"
)

var (
	// TransactionIDRegEx represents a regular expression used to validate the
	// scheme of transaction IDs. Transaction IDs provided via HTTP headers can be
	// validated using this.
	TransactionIDRegEx = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{30,62}[A-Za-z0-9]$`)
)

// IsValidTransactionID is a convenience method to validate a transaction ID.
// Internally TransactionIDRegEx is used.
func IsValidTransactionID(transactionID string) bool {
	return TransactionIDRegEx.MatchString(transactionID)
}
