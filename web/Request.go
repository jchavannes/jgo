package web

import (
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	HttpRequest http.Request
}

func (r *Request) GetCsrfToken() (string, error) {
	if r.HttpRequest.Method != "POST" {
		return "", errors.New("Not a POST request.")
	}
	csrfToken := r.HttpRequest.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		return "", errors.New("Header empty or not set.")
	}
	return csrfToken, nil
}

func (r *Request) GetFormValue(key string) string {
	r.HttpRequest.ParseForm()
	return r.HttpRequest.Form.Get(key)
}

func (r *Request) GetFormValueBool(key string) bool {
	r.HttpRequest.ParseForm()
	return getBoolFromString(strings.ToLower(r.HttpRequest.Form.Get(key)))
}

func getBoolFromString(val string) bool {
	if val == "true" {
		return true
	} else if val == "false" {
		return false
	}
	i, _ := strconv.Atoi(val)
	return i == 1
}

func (r *Request) GetFormValueInt(key string) int {
	r.HttpRequest.ParseForm()
	valString := r.HttpRequest.Form.Get(key)
	i, _ := strconv.Atoi(valString)
	return i
}

func (r *Request) GetFormValueInt64(key string) int64 {
	r.HttpRequest.ParseForm()
	valString := r.HttpRequest.Form.Get(key)
	i, _ := strconv.ParseInt(valString, 0, 64)
	return i
}

func (r *Request) GetFormValueUint(key string) uint {
	return uint(r.GetFormValueInt(key))
}

func (r *Request) GetFormValueUint64(key string) uint64 {
	return uint64(r.GetFormValueInt64(key))
}

func (r *Request) GetFormValueFloat(key string) float32 {
	r.HttpRequest.ParseForm()
	valString := r.HttpRequest.Form.Get(key)
	i, _ := strconv.ParseFloat(valString, 32)
	return float32(i)
}

func (r *Request) GetFormValueFloat64(key string) float64 {
	r.HttpRequest.ParseForm()
	valString := r.HttpRequest.Form.Get(key)
	i, _ := strconv.ParseFloat(valString, 64)
	return i
}

func (r *Request) GetFormValueSlice(key string) []string {
	r.HttpRequest.ParseForm()
	value, ok := r.HttpRequest.Form[key+"[]"]
	if !ok {
		return []string{}
	}
	return value
}

func (r *Request) GetFormValueIntSlice(key string) []int {
	r.HttpRequest.ParseForm()
	value, ok := r.HttpRequest.Form[key+"[]"]
	if !ok {
		return []int{}
	}
	var ints []int
	for _, item := range value {
		i, _ := strconv.Atoi(item)
		ints = append(ints, i)
	}
	return ints
}

func (r *Request) GetUrlNamedQueryVariable(key string) string {
	vars := mux.Vars(&r.HttpRequest)
	return vars[key]
}

func (r *Request) GetUrlNamedQueryVariableInt(key string) int {
	i, _ := strconv.Atoi(r.GetUrlNamedQueryVariable(key))
	return i
}

func (r *Request) GetUrlNamedQueryVariableUInt(key string) uint {
	return uint(r.GetUrlNamedQueryVariableInt(key))
}

func (r *Request) GetUrlParameter(key string) string {
	return r.HttpRequest.URL.Query().Get(key)
}

func (r *Request) GetUrlParameterSet(key string) bool {
	_, ok := r.HttpRequest.URL.Query()[key]
	return ok
}

func (r *Request) GetUrlParameterInt(key string) int {
	i, _ := strconv.Atoi(r.GetUrlParameter(key))
	return i
}

func (r *Request) GetUrlParameterBool(key string) bool {
	return getBoolFromString(strings.ToLower(r.GetUrlParameter(key)))
}

func (r *Request) GetUrlParameterUInt(key string) uint {
	return uint(r.GetUrlParameterInt(key))
}

func (r *Request) GetHeader(key string) string {
	return r.HttpRequest.Header.Get(key)
}

func (r *Request) GetCookie(key string) string {
	cookie, _ := r.HttpRequest.Cookie(key)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

func (r *Request) GetCookieBool(key string) bool {
	return getBoolFromString(strings.ToLower(r.GetCookie(key)))
}

func (r *Request) GetURI() string {
	return r.HttpRequest.RequestURI
}

func (r *Request) GetBody() []byte {
	body, _ := ioutil.ReadAll(r.HttpRequest.Body)
	return body
}

func (r *Request) GetPotentialFilename() string {
	return r.HttpRequest.RequestURI[1:len(r.HttpRequest.RequestURI)]
}

func (r *Request) GetSourceIP() string {
	cfIp := r.GetHeader("CF-Connecting-IP")
	if cfIp != "" {
		return cfIp
	}
	forwarded := r.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	host, _, _ := net.SplitHostPort(r.HttpRequest.RemoteAddr)
	return host
}
