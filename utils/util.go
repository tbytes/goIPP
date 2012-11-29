package utils

import (
	"time"
)

type counter struct {
	startTime, endTime, stopTime time.Time
	end, current, by             int
	finished                     bool
	fc							 func(a, b interface{}) interface{}	
}

func (c *counter) Finished() bool {
	if c.stopTime.After(time.Now()) {
		c.finished = true
	}
	if c.current >= c.end {
		c.finished = true
	}
	return c.finished
}

func (c *counter) plus() bool {
	c.current += c.by
	return c.Finished()
}

func (c *counter) minus() bool {
	c.current -= c.by
	return c.Finished()
}

func (c *counter) Get() int {
	return c.current
}

func (c *counter) SetEnd(t time.Time) {
	c.endTime = t
	return
}

func (c *counter) Max(i int) {
	c.end = i
	return
}

func (c *counter) GT(i int) bool {
	if c.current > i {
		return true
	}
	return false
}

func (c *counter) GTEq(i int) bool {
	if c.current >= i {
		return true
	}
	return false
}

func (c *counter) LT(i int) bool {
	if c.current < i {
		return true
	}
	return false
}

func (c *counter) LTEq(i int) bool {
	if c.current <= i {
		return true
	}
	return false
}

func (c *counter) Eq(i int) bool {
	if c.current <= i {
		return true
	}
	return false
}

func (c *counter) Lp(fc func(a, b interface{}) interface{}, i int) {
	c.fc = fc
	c.end = i
	return 
}

func (c *counter) Plus() bool {
	return c.plus()
}

func (c *counter) Minus() bool {
	return c.plus()
}

func Start(i int) (c counter) {
	c.by = i
	c.current = 0
	c.finished = false
	c.startTime = time.Now()
	return
}
