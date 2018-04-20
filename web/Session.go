package web

import (
	"crypto/rand"
	"encoding/hex"
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
		s.SetCsrfToken(CreateToken())
	}
	return csrfTokens[s.CookieId]
}

func (s *Session) SetCsrfToken(csrfToken string) {
	csrfTokens[s.CookieId] = csrfToken
}

func CreateToken() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
