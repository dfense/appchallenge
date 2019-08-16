package service

import (
	"errors"
	"sync"
	"time"

	"github.com/dfense/appchallenge"
	"github.com/dfense/appchallenge/client"
	log "github.com/sirupsen/logrus"
)

var (
	ErrQueueEmpty       = errors.New("queue empty") // return error if queue us ever empty
	ErrNonamesRetrieved = errors.New("no names retrieved from service")
	lowMarker           = 10       // when size of cache fall below this, start routine to add more
	defaultBatchSize    = 100      // fetch names in sizes of
	onceNameCache       sync.Once  // one time type
	instance            *nameCache // package instance constructor
)

// namceCache - keeps a supply of names, to avoid going to service as often
// service also has a rate limiter, but retrieving batches is more lienent
type nameCache struct {
	sync.Mutex // build lock into struct
	people     []appchallenge.Person
	updating   bool // actively getting more data from ext service
}

// GetNameCache - singleton for name cache.
func GetNameCache() *nameCache {
	onceNameCache.Do(func() {
		instance = &nameCache{}
	})
	return instance
}

// enqueue - add more people to the cache, called by update.
// do not call directly. no locking
func (nc *nameCache) enqueue(people []appchallenge.Person) {
	// mutex lock
	nc.people = append(nc.people, people...) // enqueue variadic
}

// Dequeue - select the first element from the slice and then
// delete it
func (nc *nameCache) Dequeue() (*appchallenge.Person, error) {
	// check size, and get more if low
	// mutex lock
	nc.Lock()
	defer nc.Unlock()

	// restock the queue once we fall below watermark
	// run update in the background, and keep it ahead of dequeue going empty
	if len(nc.people) < lowMarker {
		go nc.update()
	}

	// check to make sure we aren't empty and can't rebuild
	attempts := 0
	for len(nc.people) == 0 {
		if attempts > 50 {
			return nil, ErrQueueEmpty
		}
		attempts++
		time.Sleep(time.Millisecond * 100)
	}

	// make a copy then delete original
	tmpCopy := nc.people[0]   // first element
	nc.people = nc.people[1:] // dequeue
	return &tmpCopy, nil
}

// internal method to update more cache of people
// lock should have been obtained from dequeue. do not call direct
func (nc *nameCache) update() error {

	log.Debug("updating cache")
	nc.updating = true
	names, err := client.GetNames(defaultBatchSize)
	if err != nil {
		log.Error(err)
		nc.updating = false
		return err
	}

	if len(names) == 0 {
		nc.updating = false
		return ErrNonamesRetrieved
	}
	nc.enqueue(names)

	nc.updating = false
	log.Debugf("cache size now %d", len(nc.people))
	return nil
}
