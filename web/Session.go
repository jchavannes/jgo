package web

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Session struct {
	CookieId string
}

var csrfTokens map[string]string

func (s *Session) GetCsrfToken() string {
	if csrfTokens == nil {
		csrfTokens = make(map[string]string)
	}
	if _, ok := csrfTokens[s.CookieId]; !ok {
		csrfTokens[s.CookieId] = CreateToken()
	}
	return csrfTokens[s.CookieId]
}

func CreateToken() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
