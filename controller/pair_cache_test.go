package controller

import "testing"

func Test_pairCache(t *testing.T) {
	testCases := []struct {
		name   string
		size   int
		input  []stringPair
		hits   []stringPair
		misses []stringPair
	}{
		{
			name: "case 0: basic",
			size: 3,
			input: []stringPair{
				{A: "a", B: ""},
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			hits: []stringPair{
				{A: "a", B: ""},
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			misses: []stringPair{
				{A: "x", B: ""},
				{A: "", B: "y"},
				{A: "x", B: "y"},
			},
		},
		{
			name: "case 1: overflow (size is 1, but 3 elements put)",
			size: 2,
			input: []stringPair{
				{A: "x", B: "y"},
				{A: "a", B: ""},
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			hits: []stringPair{
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			misses: []stringPair{
				{A: "x", B: "y"},
				{A: "a", B: ""},
			},
		},
		{
			name: "case 2: underflow (size is 10, but 3 elements are put)",
			size: 10,
			input: []stringPair{
				{A: "a", B: ""},
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			hits: []stringPair{
				{A: "a", B: ""},
				{A: "", B: "b"},
				{A: "a", B: "b"},
			},
			misses: []stringPair{
				{A: "x", B: ""},
				{A: "", B: "y"},
				{A: "x", B: "y"},
			},
		},
		{
			name: "case 3: empty string hit",
			size: 2,
			input: []stringPair{
				{A: "", B: ""},
				{A: "a", B: "b"},
			},
			hits: []stringPair{
				{A: "", B: ""},
			},
		},
		{
			name: "case 4: empty string miss",
			size: 1,
			input: []stringPair{
				{A: "a", B: "b"},
			},
			misses: []stringPair{
				{A: "", B: ""},
			},
		},
		{
			name: "case 5: same input do not cause overflow",
			size: 2,
			input: []stringPair{
				{A: "a", B: "b"},
				{A: "a", B: "b"},
				{A: "a", B: "b"},
				{A: "xxx", B: "yyy"},
				{A: "a", B: "b"},
				{A: "a", B: "b"},
				{A: "a", B: "b"},
				{A: "a", B: "b"},
				{A: "a", B: "b"},
			},
			hits: []stringPair{
				{A: "a", B: "b"},
				{A: "xxx", B: "yyy"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := newPairCache(tc.size)

			for _, p := range tc.input {
				cache.Put(p.A, p.B)
			}

			for _, p := range tc.hits {
				contains := cache.Contains(p.A, p.B)
				if !contains {
					t.Errorf("cache.Contains(%q, %q) == %v, want %v", p.A, p.B, contains, true)
				}
			}

			for _, p := range tc.misses {
				contains := cache.Contains(p.A, p.B)
				if contains {
					t.Errorf("cache.Contains(%q, %q) == %v, want %v", p.A, p.B, contains, false)
				}
			}
		})
	}
}
