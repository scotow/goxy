package server

import (
	"errors"
	. "github.com/scotow/goxy/common"
	"sync"
)

const (
	maxAttempts = 20
)

var (
	errTooManyCreationAttempts = errors.New("too many attempts to find a free id")
	errConnectionNotFound      = errors.New("cannot find connection matching this token")
)

func newConnMapping() *connMapping {
	return &connMapping{
		mapping: make(map[string]*Conn),
	}
}

type connMapping struct {
	mapping map[string]*Conn
	lock    sync.RWMutex
}

func (c *connMapping) addConn(conn *Conn) error {
	var id *Id
	attempts := 0
	c.lock.RLock()

	for {
		id = NewRandomId()

		if c.mapping[id.Token()] == nil {
			break
		}

		attempts += 1
		if attempts == maxAttempts {
			c.lock.RUnlock()
			return errTooManyCreationAttempts
		}
	}

	c.lock.RUnlock()

	c.lock.Lock()
	c.mapping[id.Token()] = conn
	c.lock.Unlock()

	conn.id = id
	return nil
}

func (c *connMapping) getConn(token string) (*Conn, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	conn := c.mapping[token]
	if conn == nil {
		return nil, errConnectionNotFound
	}

	return c.mapping[token], nil
}
