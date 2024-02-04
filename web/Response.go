package web

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response objects are passed to handlers to respond to requests.
// Includes abstracted access to request and session information
type Response struct {
	Helper  map[string]interface{}
	funcMap map[string]interface{}
	StartTs time.Time
	Pattern string
	Request Request
	Server  *Server
	Session Session
	Static  bool
	code    int
	Writer  http.ResponseWriter
}

func (r *Response) SetFuncMap(funcMap map[string]interface{}) {
	r.funcMap = template.FuncMap(funcMap)
}

// Checks that CSRF token in request matches one for session.
// Tokens are kept in memory and do not persist between instances or restarts.
func (r *Response) IsValidCsrf() bool {
	requestCsrfToken, err := r.Request.GetCsrfToken()
	return err == nil && requestCsrfToken != "" && requestCsrfToken == r.GetCsrfToken()
}

func (r *Response) GetCsrfToken() string {
	if r.Server.GetCsrfToken != nil {
		return r.Server.GetCsrfToken(r.Session.CookieId)
	}
	return r.Session.GetCsrfToken()
}

func (r *Response) SetCookie(key string, value string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(value),
		Path:     "/",
		HttpOnly: true,
		Secure:   !r.Server.InsecureCookie,
		MaxAge:   int(time.Hour) * 24 * 30,
	}
	http.SetCookie(r.Writer, &cookie)
}

// Sets a new session cookie.
func (r *Response) ResetOrCreateSession() {
	r.SetSessionCookie(CreateToken())
}

func (r *Response) SetSessionCookie(cookie string) {
	r.Session = Session{
		CookieId: cookie,
	}
	r.SetCookie(r.Server.GetCookieName(), r.Session.CookieId)
}

// Either gets existing session token or creates a new one.
func (r *Response) InitSession() {
	cookie := r.Request.GetCookie(r.Server.GetCookieName())
	cookieByte, err := hex.DecodeString(cookie)
	if err == nil && cookie != "" && len(cookieByte) == 32 {
		r.Session = Session{
			CookieId: cookie,
		}
	} else {
		r.ResetOrCreateSession()
	}
}

func (r *Response) ResponseCodeSet() bool {
	return r.code != 0
}

func (r *Response) GetResponseCode() int {
	responseCode := r.code
	if responseCode == 0 {
		responseCode = http.StatusOK
	}
	return responseCode
}

func (r *Response) SetResponseCode(code int) error {
	if r.ResponseCodeSet() {
		return errors.New("response code already set")
	}
	r.Writer.WriteHeader(code)
	r.code = code
	return nil
}

func (r *Response) SetHeader(key string, value string) {
	r.Writer.Header().Set(key, value)
}

func (r *Response) Write(s string) {
	r.Writer.Write([]byte(s))
}

func (r *Response) WriteJson(i interface{}, pretty bool) {
	var text []byte
	if pretty {
		text, _ = json.MarshalIndent(i, "", "  ")
	} else {
		text, _ = json.Marshal(i)
	}
	r.SetHeader("Content-Type", "application/json")
	r.Write(string(text))
}

func (r *Response) Render() error {
	requestURI := r.Request.GetURI()
	templateName := strings.Split(requestURI[1:], "?")[0]
	return r.RenderTemplate(templateName)
}

func (r *Response) RenderTemplate(templateName string) error {
	renderer, err := GetRenderer(r.Server.TemplatesDir)
	if err != nil {
		jerr.Get("error getting renderer", err).Print()
	}

	if r.funcMap != nil {
		renderer.SetFuncMap(r.funcMap)
	}

	r.SetHeader("Content-Type", "text/html")

	if strings.HasPrefix(templateName, "/") {
		templateName = templateName[1:]
	}

	var indexTemplate string
	if templateName == "" {
		indexTemplate = "index.html"
	} else {
		indexTemplate = templateName + "/index.html"
	}

	err = renderer.Render([]string{
		templateName + ".html",
		indexTemplate,
	}, r.Writer, r.Helper)

	if jerr.HasError(err, UnableToFindTemplateErrorMsg) {
		r.Error(jerr.Get("template not found", err), http.StatusNotFound)
		err = renderer.Render([]string{
			"404.html",
		}, r.Writer, r.Helper)
	} else if err != nil {
		r.Error(jerr.Get("error rendering template", err), http.StatusInternalServerError)
	}

	return err
}

func (r *Response) SetRedirect(location string) {
	r.SetHeader("Location", location)
	r.SetResponseCode(http.StatusFound)
}

func (r *Response) GetWebSocket() (*Socket, error) {
	upgrader := websocket.Upgrader{
		EnableCompression: true,
	}
	conn, err := upgrader.Upgrade(r.Writer, &r.Request.HttpRequest, nil)
	if err != nil {
		return nil, err
	}
	r.code = http.StatusSwitchingProtocols
	return &Socket{ws: conn}, nil
}

func (r *Response) Error(err error, responseCode int) {
	r.SetResponseCode(responseCode)
	err = jerr.Get(fmt.Sprintf("error with request: %#v", r.Request.HttpRequest.URL.Path), err)
	jlog.Log(err)
	if r.Server.ErrorHandler != nil {
		r.Server.ErrorHandler(r, err)
	}
}

const LogCompleteExtraMessageHelperKey = "LogCompleteExtraMessage"

func (r *Response) LogComplete() {
	var extraMessage string
	if val, ok := r.Helper[LogCompleteExtraMessageHelperKey]; ok {
		extraMessage = fmt.Sprintf(" (%s)", val)
	}
	jlog.Logf("%s %s %s %d%s\n",
		jutil.ShortHash(r.Session.CookieId),
		r.Request.GetSourceIP(),
		r.Request.HttpRequest.URL.Path,
		r.GetResponseCode(),
		extraMessage,
	)
}
