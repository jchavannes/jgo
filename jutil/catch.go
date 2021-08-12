package jutil

type Catch struct {
	Panic interface{}
}

func (c *Catch) Try(f func()) {
	defer func() {
		c.Panic = recover()
	}()
	f()
}
