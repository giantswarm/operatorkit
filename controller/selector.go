package controller

import (
	"k8s.io/apimachinery/pkg/labels"
)

type Selector interface {
	Matches(labels Labels) bool
}

type Labels = labels.Labels

type internalLabels map[string]string

func (l internalLabels) Has(label string) bool {
	_, ok := l[label]
	return ok
}

func (l internalLabels) Get(label string) string {
	return l[label]
}
