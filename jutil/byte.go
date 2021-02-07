package jutil

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
)

func ByteReverse(d []byte) []byte {
	l := len(d)
	var n = make([]byte, l)
	for i := 0; i < l; i++ {
		n[i] = d[l-1-i]
	}
	return n
}

func RemoveDupesAndEmpties(data [][]byte) [][]byte {
	for i := 0; i < len(data); i++ {
		if len(data[i]) == 0 {
			data = append(data[:i], data[i+1:]...)
			i--
			continue
		}
		for g := i + 1; g < len(data); g++ {
			if bytes.Equal(data[i], data[g]) {
				data = append(data[:g], data[g+1:]...)
				g--
			}
		}
	}
	return data
}

func RemoveDupes(uints []uint) []uint {
	for i := 0; i < len(uints); i++ {
		for g := i + 1; g < len(uints); g++ {
			if uints[i] == uints[g] {
				uints = append(uints[:g], uints[g+1:]...)
				g--
			}
		}
	}
	return uints
}

func InByteArray(needle []byte, haystack [][]byte) bool {
	for _, hay := range haystack {
		if bytes.Equal(needle, hay) {
			return true
		}
	}
	return false
}

func GetUint64(data []byte) uint64 {
	const size = 8
	tmp := make([]byte, size)
	if len(data) > size {
		tmp = data[len(data)-size:]
	} else {
		copy(tmp[size-len(data):], data)
	}
	return binary.BigEndian.Uint64(tmp)
}

func GetUint(data []byte) uint {
	return uint(GetUint32(data))
}

func GetUint32(data []byte) uint32 {
	const size = 4
	tmp := make([]byte, size)
	if len(data) > size {
		tmp = data[len(data)-size:]
	} else {
		copy(tmp[size-len(data):], data)
	}
	return binary.LittleEndian.Uint32(tmp)
}

func GetUint32Data(i uint32) []byte {
	var b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

func GetInt(data []byte) int {
	return int(GetUint32(data))
}

func GetIntData(i int) []byte {
	return GetUint32Data(uint32(i))
}

func GetInt32(data []byte) int32 {
	return int32(GetUint32(data))
}

func GetInt32Data(i int32) []byte {
	return GetUint32Data(uint32(i))
}

func GetInt64Data(i int64) []byte {
	var b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func GetInt64(d []byte) int64 {
	i := binary.LittleEndian.Uint64(d)
	return int64(i)
}

func GetInt64DataBig(i int64) []byte {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func GetInt64Big(d []byte) int64 {
	i := binary.BigEndian.Uint64(d)
	return int64(i)
}

func HasPrefix(b []byte, prefix []byte) bool {
	if len(prefix) > len(b) {
		return false
	}
	return bytes.Equal(b[:len(prefix)], prefix)
}

// https://stackoverflow.com/a/40678026/744298
func CombineBytes(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	b := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(b[i:], s)
	}
	return b
}

// https://stackoverflow.com/a/45506459/744298
func AllZeros(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func BytePad(b []byte, size int) []byte {
	if len(b) < size {
		var eb = make([]byte, size-len(b))
		b = append(b, eb...)
	}
	return b
}

func GetSha256Hash(script []byte) []byte {
	s := sha256.Sum256(script)
	return s[:]
}

func GetByteMd5Int(b []byte) uint {
	s := md5.Sum(b)
	i := binary.BigEndian.Uint32(s[:8])
	return uint(i)
}
