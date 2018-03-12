package web

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

type UrlParamType string

const UrlParamInteger UrlParamType = "integer"

type UrlParam struct {
	Id   string
	Type UrlParamType
}

func (u UrlParam) UrlPart() string {
	switch u.Type {
	case UrlParamInteger:
		return "{" + u.Id + ":[0-9]+}"
	}
	fmt.Println(jerr.New("unknown url param type").Error())
	return ""
}
