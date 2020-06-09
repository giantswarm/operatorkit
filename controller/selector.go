package controller

import (
	"k8s.io/apimachinery/pkg/labels"
)

type Labels = labels.Labels

type Selector interface {
	Matches(labels Labels) bool
}

func NewSelector(matchesFunc func(labels Labels) bool) Selector {
	return &internalSelector{
		matchesFunc: matchesFunc,
	}
}

type internalSelector struct {
	matchesFunc func(labels Labels) bool
}

func (s *internalSelector) Matches(labels Labels) bool {
	return s.matchesFunc(labels)
}

type internalLabels map[string]string

func (l internalLabels) Has(label string) bool {
	_, ok := l[label]
	return ok
}

func (l internalLabels) Get(label string) string {
	return l[label]
}
