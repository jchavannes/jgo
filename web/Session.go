package web

import (
	"io"
	"encoding/base64"
	"crypto/rand"
)

type Session struct {
	CookieId string
}

func CreateSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
