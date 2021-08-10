package jutil

type Catch struct {
	Panic interface{}
}

func (r *Catch) Try(f func()) {
	defer func() {
		r.Panic = recover()
	}()
	f()
}
