package controller

import "sync"

// pairCache is a very simple last in first out in-memory cache storing string
// pairs.
type pairCache struct {
	data      []stringPair
	dataMutex *sync.Mutex
}

func newPairCache(size int) *pairCache {
	if size < 1 {
		// Yes this panics. But if size is less than 0 it will panic
		// anyway in slice creation.
		panic("size must be bigger than 0")
	}

	p := &pairCache{
		data:      make([]stringPair, 0, size),
		dataMutex: &sync.Mutex{},
	}

	return p
}

// Put inserts and item into the cache.
func (p *pairCache) Put(a, b string) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()

	if p.containsThreadUnsafe(a, b) {
		return
	}

	pair := stringPair{
		A: a,
		B: b,
	}

	if len(p.data) != cap(p.data) {
		p.data = append(p.data, pair)
		return
	}

	copy(p.data, p.data[1:])
	p.data[len(p.data)-1] = pair
}

func (p *pairCache) Contains(a, b string) bool {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()

	return p.containsThreadUnsafe(a, b)
}

func (p *pairCache) containsThreadUnsafe(a, b string) bool {
	// Assuming we search for most recent items more often searching from
	// the end will return faster if this is a hit.
	for i := len(p.data) - 1; i >= 0; i-- {
		pair := p.data[i]
		if pair.A == a && pair.B == b {
			return true
		}
	}

	return false
}

type stringPair struct {
	A string
	B string
}
