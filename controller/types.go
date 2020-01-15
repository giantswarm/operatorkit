package controller

// mapLabels implements labels.Labels. See
// https://godoc.org/k8s.io/apimachinery/pkg/labels#Labels.
type mapLabels struct {
	m map[string]string
}

func newMapLabels(m map[string]string) *mapLabels {
	return &mapLabels{
		m: m,
	}
}

func (l *mapLabels) Has(label string) bool {
	_, ok := l.m[label]
	return ok
}

func (l *mapLabels) Get(label string) string {
	v, _ := l.m[label]
	return v
}
