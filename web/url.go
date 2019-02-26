package web

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

type UrlParamType string

const (
	UrlParamInteger      UrlParamType = "integer"
	UrlParamString       UrlParamType = "string"
	UrlParamString2      UrlParamType = "string2"
	UrlParamAlphaNumeric UrlParamType = "alpha-numeric"
	UrlParamAny          UrlParamType = "any"
)

type UrlParam struct {
	Id   string
	Type UrlParamType
}

func (u UrlParam) UrlPart() string {
	switch u.Type {
	case UrlParamInteger:
		return "{" + u.Id + ":[0-9]+}"
	case UrlParamString, UrlParamAlphaNumeric:
		return "{" + u.Id + ":[A-Za-z0-9]+}"
	case UrlParamString2:
		return "{" + u.Id + ":[A-Za-z0-9-_]+}"
	case UrlParamAny:
		return "{" + u.Id + ":.*}"
	}
	fmt.Println(jerr.New("unknown url param type").Error())
	return ""
}
