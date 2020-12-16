package selector

import (
	"testing"

	"k8s.io/apimachinery/pkg/labels"
)

func Test_Selector_compatibility(*testing.T) {
	var _ Selector = labels.NewSelector()
}
