package controller

import "testing"

func Test_pairCache(t *testing.T) {
	testCases := []struct {
		name   string
		size   int
		input  []string
		hits   []string
		misses []string
	}{
		{
			name:   "case 0: basic",
			size:   4,
			input:  []string{"a", "b", "", " "},
			hits:   []string{"a", "b", "", " "},
			misses: []string{"x", "y", "        ", "c"},
		},
		{
			name:   "case 1: overflow (size is 2, but 4 elements put)",
			size:   2,
			input:  []string{"a", "b", "c", "d"},
			hits:   []string{"c", "d"},
			misses: []string{"a", "b"},
		},
		{
			name:   "case 2: underflow (size is 10, but 3 elements are put)",
			size:   10,
			input:  []string{"a", "b", "c"},
			hits:   []string{"a", "b", "c"},
			misses: []string{"x", "y"},
		},
		{
			name:  "case 3: empty string hit",
			size:  2,
			input: []string{" ", "   "},
			hits:  []string{" ", "   "},
		},
		{
			name:   "case 4: empty string miss",
			size:   2,
			input:  []string{"a", "b"},
			misses: []string{" ", "   "},
		},
		{
			name:  "case 5: same input do not cause overflow",
			size:  2,
			input: []string{"a", "a", "b", "a", "a"},
			hits:  []string{"a", "b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := newFifoCache(tc.size)

			for _, s := range tc.input {
				cache.Put(s)
			}

			for _, s := range tc.hits {
				contains := cache.Contains(s)
				if !contains {
					t.Errorf("cache.Contains(%q) == %v, want %v", s, contains, true)
				}
			}

			for _, s := range tc.misses {
				contains := cache.Contains(s)
				if contains {
					t.Errorf("cache.Contains(%q) == %v, want %v", s, contains, false)
				}
			}
		})
	}
}
