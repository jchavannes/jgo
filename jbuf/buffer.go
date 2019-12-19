package jbuf

type Item interface {
}

type buffer struct {
	new   chan interface{}
	done  chan struct{}
	buf   []interface{}
	inUse bool

	handle func([]interface{})
}

func NewBuffer(handle func([]interface{})) *buffer {
	var updater = &buffer{
		handle: handle,
	}
	updater.init()
	return updater
}

func (u *buffer) init() {
	u.new = make(chan interface{})
	u.done = make(chan struct{})
	go func() {
		for {
			select {
			case n := <-u.new:
				if u.inUse {
					u.buf = append(u.buf, n)
				} else {
					u.inUse = true
					go u.process([]interface{}{n})
				}
			case <-u.done:
				if len(u.buf) > 0 {
					go u.process(u.buf)
					u.buf = []interface{}{}
				} else {
					u.inUse = false
				}
			}
		}
	}()
}

func (u *buffer) Buffer(item interface{}) {
	u.new <- item
}

func (u *buffer) process(buf []interface{}) {
	u.handle(buf)
	u.done <- struct{}{}
}
