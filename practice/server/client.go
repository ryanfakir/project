package server

import (
	"context"
	"sync"
	"time"
)

// Entry in map
type Entry struct {
	res   Result
	ready chan struct{}
}

// Client with cache
type Client struct {
	DataServiceImpl DataService
	mu              sync.Mutex
	cache           map[string]*Entry
	history         map[string]*Entry
}

// CallService will do caching
func (c *Client) CallService(ctx context.Context, query string) Result {
	c.mu.Lock()
	e := c.cache[query]
	// upload e to history cache
	c.history[query] = e
	t1 := time.Now()
	if e == nil || t1.After(e.res.QueryTime.Add(time.Minute)) {
		e = &Entry{ready: make(chan struct{})}
		c.cache[query] = e
		c.mu.Unlock()
		e.res, e.res.Err = c.DataServiceImpl.CallService(ctx, query)
		close(e.ready)
	} else {
		c.mu.Unlock()
		<-e.ready
	}
	return e.res
}

// NewClient Pointer
func NewClient() *Client {
	return &Client{
		DataServiceImpl: NewService(),
		cache:           make(map[string]*Entry),
		history:         make(map[string]*Entry),
	}
}
