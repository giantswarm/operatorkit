package controller

type Object interface {
	GetAnnotations() map[string]string
	GetLabels() map[string]string
}

type Selector interface {
	Matches(obj Object) bool
}

func NewSelector(matchesFunc func(obj Object) bool) Selector {
	return &internalSelector{
		matchesFunc: matchesFunc,
	}
}

func NewSelectorEverything() Selector {
	return NewSelector(func(obj Object) bool {
		return true
	})
}

type internalSelector struct {
	matchesFunc func(obj Object) bool
}

func (s *internalSelector) Matches(obj Object) bool {
	return s.matchesFunc(obj)
}
