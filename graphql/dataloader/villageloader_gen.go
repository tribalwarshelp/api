// Code generated by github.com/vektah/dataloaden, DO NOT EDIT.

package dataloader

import (
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"sync"
	"time"
)

// VillageLoaderConfig captures the config to create a new VillageLoader
type VillageLoaderConfig struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(keys []int) ([]*twmodel.Village, []error)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewVillageLoader creates a new VillageLoader given a fetch, wait, and maxBatch
func NewVillageLoader(config VillageLoaderConfig) *VillageLoader {
	return &VillageLoader{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// VillageLoader batches and caches requests
type VillageLoader struct {
	// this method provides the data for the loader
	fetch func(keys []int) ([]*twmodel.Village, []error)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
	cache map[int]*twmodel.Village

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *villageLoaderBatch

	// mutex to prevent races
	mu sync.Mutex
}

type villageLoaderBatch struct {
	keys    []int
	data    []*twmodel.Village
	error   []error
	closing bool
	done    chan struct{}
}

// Load a Village by key, batching and caching will be applied automatically
func (l *VillageLoader) Load(key int) (*twmodel.Village, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a Village.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *VillageLoader) LoadThunk(key int) func() (*twmodel.Village, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (*twmodel.Village, error) {
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &villageLoaderBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() (*twmodel.Village, error) {
		<-batch.done

		var data *twmodel.Village
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		if err == nil {
			l.mu.Lock()
			l.unsafeSet(key, data)
			l.mu.Unlock()
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *VillageLoader) LoadAll(keys []int) ([]*twmodel.Village, []error) {
	results := make([]func() (*twmodel.Village, error), len(keys))

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

	villages := make([]*twmodel.Village, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range results {
		villages[i], errors[i] = thunk()
	}
	return villages, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Villages.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *VillageLoader) LoadAllThunk(keys []int) func() ([]*twmodel.Village, []error) {
	results := make([]func() (*twmodel.Village, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([]*twmodel.Village, []error) {
		villages := make([]*twmodel.Village, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range results {
			villages[i], errors[i] = thunk()
		}
		return villages, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
func (l *VillageLoader) Prime(key int, value *twmodel.Village) bool {
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
		cpy := *value
		l.unsafeSet(key, &cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *VillageLoader) Clear(key int) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

func (l *VillageLoader) unsafeSet(key int, value *twmodel.Village) {
	if l.cache == nil {
		l.cache = map[int]*twmodel.Village{}
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *villageLoaderBatch) keyIndex(l *VillageLoader, key int) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *villageLoaderBatch) startTimer(l *VillageLoader) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *villageLoaderBatch) end(l *VillageLoader) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
