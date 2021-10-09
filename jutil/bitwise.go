package jutil

type Bitwise []byte

func (b *Bitwise) fill(index int) {
	if len(*b) <= index {
		for i := len(*b); i <= index; i++ {
			*b = append(*b, 0)
		}
	}
}

func (b *Bitwise) Set(index int, flag byte) {
	b.fill(index)
	(*b)[index] = (*b)[index] | flag
}

func (b *Bitwise) Clear(index int, flag byte) {
	b.fill(index)
	(*b)[index] = (*b)[index] &^ flag
}

func (b *Bitwise) Toggle(index int, flag byte) {
	b.fill(index)
	(*b)[index] = (*b)[index] ^ flag
}

func (b Bitwise) Has(index int, flag byte) bool {
	if len(b) <= index {
		return false
	}
	return b[index]&flag != 0
}
