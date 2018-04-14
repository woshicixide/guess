package vote

import (
	// "log"
	"sync"
)

// 存储单位
type container struct {
	votes map[string]*codeSto
	rw    *sync.RWMutex
}

func (self *container) get(code string) *codeSto {
	var c *codeSto

	self.rw.RLock()
	c, exists := self.votes[code]
	self.rw.RUnlock()
	// log.Fatal(self)

	if !exists {
		c = New()
		self.set(code, c)
	}
	return c
}

func (self *container) set(code string, c *codeSto) {
	self.rw.Lock()
	defer self.rw.Unlock()

	self.votes[code] = c
}

func (self *container) reset() {
	self.rw.Lock()
	defer self.rw.Unlock()
	self.votes = nil
	self.votes = make(map[string]*codeSto)
}

func Reset() {
	con.reset()
}

var con = &container{make(map[string]*codeSto), new(sync.RWMutex)}
