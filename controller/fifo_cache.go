package controller

import "sync"

// fifoCache is a very simple first in first out in-memory cache storing string
// pairs.
type fifoCache struct {
	data      []string
	dataMutex *sync.Mutex
}

func newFifoCache(size int) *fifoCache {
	if size < 1 {
		// Yes this panics. But if size is less than 0 it will panic
		// anyway in slice creation.
		panic("size must be bigger than 0")
	}

	f := &fifoCache{
		data:      make([]string, 0, size),
		dataMutex: &sync.Mutex{},
	}

	return f
}

// Put inserts and item into the cache.
func (f *fifoCache) Put(s string) {
	f.dataMutex.Lock()
	defer f.dataMutex.Unlock()

	if f.containsThreadUnsafe(s) {
		return
	}

	if len(f.data) != cap(f.data) {
		f.data = append(f.data, s)
		return
	}

	copy(f.data, f.data[1:])
	f.data[len(f.data)-1] = s
}

func (f *fifoCache) Contains(s string) bool {
	f.dataMutex.Lock()
	defer f.dataMutex.Unlock()

	return f.containsThreadUnsafe(s)
}

func (f *fifoCache) containsThreadUnsafe(s string) bool {
	// Assuming we search for most recent items more often searching from
	// the end will return faster if this is a hit.
	for i := len(f.data) - 1; i >= 0; i-- {
		if s == f.data[i] {
			return true
		}
	}

	return false
}
