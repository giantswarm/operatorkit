package certificatetpr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidComponent(t *testing.T) {
	tests := []struct {
		name       string
		components []ClusterComponent
		el         ClusterComponent
		res        bool
	}{
		{
			name: "el is present",
			components: []ClusterComponent{
				ClusterComponent("foobar"),
			},
			el:  ClusterComponent("foobar"),
			res: true,
		},
		{
			name: "el is not present",
			components: []ClusterComponent{
				ClusterComponent("foo"),
			},
			el:  ClusterComponent("bar"),
			res: false,
		},
		{
			name:       "components is empty",
			components: []ClusterComponent{},
			el:         ClusterComponent("foobar"),
			res:        false,
		},
	}

	for _, tc := range tests {
		res := ValidComponent(tc.el, tc.components)

		assert.Equal(t, tc.res, res)
	}
}
