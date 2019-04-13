package api

import "bytes"

// Pool is the struct that represents our buffer pool.
type Pool struct {
	c chan *bytes.Buffer
}

// New returns a new buffer pool, arbitrarily sized at 2048 bytes.
func New() *Pool {
	return &Pool{
		c: make(chan *bytes.Buffer, 2048),
	}
}

// Get gets a buffer from the pool, creating a new one if none are available.
func (p *Pool) Get() *bytes.Buffer {
	select {
	case buf := <-p.c:
		// Re-use this buffer
		return buf
	default:
		// Create a new buffer
		return &bytes.Buffer{}
	}
}

// Put returns a buffer to the pool.
func (p *Pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	select {
	case p.c <- buf:
		// Return to pool
	default:
		// Pool is full, discard buffer
	}
}
