package service

import (
	"testing"
)

// TestCache - call dequeue, and expect it to update the cache.
// call it enough for it to hit the watermark and retrieve a new batch
func TestCache(t *testing.T) {

	defaultBatchSize = 3
	lowMarker = 2

	nameCache := GetNameCache()
	person, err := nameCache.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if person == nil {
		t.Error("dequeue returned nil person")
	}

	if len(instance.people) != defaultBatchSize-1 {
		t.Errorf("expected [%d] and have [%d]", defaultBatchSize-1, len(instance.people))
	}

	person, err = nameCache.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if len(instance.people) != defaultBatchSize-2 {
		t.Errorf("expected [%d] and have [%d]", defaultBatchSize-2, len(instance.people))
	}

	// this last one should force the update to cache
	person, err = nameCache.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if person == nil {
		t.Error("cache manager did not refill at lowWatermark")
	}

	// this last would use the first new cache item
	person, err = nameCache.Dequeue()
	if err != nil {
		t.Error(err)
	}
}
