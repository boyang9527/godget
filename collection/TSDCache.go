package collection

import (
	"fmt"
	"sync"
)

type TSD interface {
	Timestamp() int64
}

type TSDCache struct {
	lock     *sync.RWMutex
	data     []TSD
	capacity int
	cursor   int
}

func NewTSDCache(capacity int) *TSDCache {
	if capacity <= 0 {
		panic("invalid TSDCache capacity")
	}
	return &TSDCache{
		lock:     &sync.RWMutex{},
		data:     make([]TSD, capacity, capacity),
		capacity: capacity,
		cursor:   0,
	}
}

func (c *TSDCache) isEmpty() bool {
	return c.cursor == 0 && c.data[c.cursor] == nil
}

func (c *TSDCache) binarySearch(t int64) int {
	if c.isEmpty() {
		return 0
	}
	var l, r int
	if c.data[c.cursor] == nil {
		l = 0
		r = c.cursor - 1
	} else {
		l = c.cursor
		r = c.cursor - 1 + c.capacity
	}

	for {
		if l > r {
			return l
		}
		m := (l + r) / 2
		if t <= c.data[m%c.capacity].Timestamp() {
			r = m - 1
		} else {
			l = m + 1
		}
	}
}

func (c *TSDCache) Put(d TSD) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.isEmpty() || d.Timestamp() >= c.data[((c.cursor-1)+c.capacity)%c.capacity].Timestamp() {
		c.data[c.cursor] = d
		c.cursor = (c.cursor + 1) % c.capacity
		return
	}

	pos := c.binarySearch(d.Timestamp())
	if pos == c.cursor && c.data[c.cursor] != nil {
		return
	}

	end := c.cursor
	if c.data[end] != nil {
		end += c.capacity
	}
	for i := end; i > pos; i-- {
		c.data[i%c.capacity] = c.data[(i-1)%c.capacity]
	}
	c.data[pos%c.capacity] = d
	c.cursor = (c.cursor + 1) % c.capacity
}

func (c *TSDCache) String() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var head, tail int
	if c.data[c.cursor] == nil {
		head = 0
		tail = c.cursor - 1
	} else {
		head = c.cursor
		tail = c.cursor + c.capacity - 1
	}

	s := make([]TSD, tail-head+1)
	for i := 0; i <= tail-head; i++ {
		s[i] = c.data[(i+head)%c.capacity]
	}
	return fmt.Sprint(s)
}

/*
 Query returns the time series with the timestamp in [start, end)
*/
func (c *TSDCache) Query(start, end int64) []TSD {
	c.lock.RLock()
	defer c.lock.RUnlock()

	from := c.binarySearch(start)
	to := c.binarySearch(end)
	length := to - from
	result := make([]TSD, length)
	for i := 0; i < length; i++ {
		result[i] = c.data[(from+i)%c.capacity]
	}
	return result
}
