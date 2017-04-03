package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jchavannes/jgo/example/db"
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/jgo/chat"
	"golang.org/x/net/html"
	"net/http"
	"time"
)

const (
	WS_CTS_SendHeartBeat = "SendHeartBeat"
	WS_CTS_SendMessage = "SendMessage"
	WS_STC_ReceiveMessage = "ReceiveMessage"
	WS_STC_UserEnter = "UserEnter"
	WS_STC_UserExit = "UserExit"
)

var (
	defaultRoute = web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			redirectOrRender(r, "lobby")
		},
	}

	signupRoute = web.Route{
		Pattern: "/signup",
		Handler: func(r *web.Response) {
			redirectOrRender(r, "lobby")
		},
	}

	signupSubmitRoute = web.Route{
		Pattern: "/signup-submit",
		CsrfProtect: true,
		Handler: func(r *web.Response) {
			username := r.Request.GetFormValue("username")
			password := r.Request.GetFormValue("password")

			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}

			_, err = session.Signup(username, password)
			if err != nil {
				fmt.Printf("Error signing up: %s\n", err)
				r.SetResponseCode(http.StatusConflict)
				return
			}
		},
	}

	loginRoute = web.Route{
		Pattern: "/login",
		Handler: func(r *web.Response) {
			redirectOrRender(r, "lobby")
		},
	}

	loginSubmitRoute = web.Route{
		Pattern: "/login-submit",
		CsrfProtect: true,
		Handler: func(r *web.Response) {
			username := r.Request.GetFormValue("username")
			password := r.Request.GetFormValue("password")

			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}

			_, err = session.Login(username, password)
			if err != nil {
				fmt.Printf("Error logging in: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.Write(err.Error())
			}
		},
	}

	logoutRoute = web.Route{
		Pattern: "/logout",
		Handler: func(r *web.Response) {
			r.ResetOrCreateSession()
			r.SetRedirect(getBaseUrl(r.Request))
		},
	}

	lobbyRoute = web.Route{
		Pattern: "/lobby",
		Handler: func(r *web.Response) {
			session, err := db.GetSession(r.Session.CookieId)
			if err == nil {
				user, err := db.GetUser(db.User{
					Id: session.UserId,
				})
				if err == nil {
					r.Helper["Username"] = user.Username
					r.Render()
					return
				}
			}
			fmt.Printf("Error getting user: %s\n", err)
			r.SetRedirect(getBaseUrl(r.Request))
		},
	}

	chatRoute = web.Route{
		Pattern: "/chat",
		Handler: func(r *web.Response) {
			conn, err := r.GetWebSocket()
			if err != nil {
				fmt.Println(err)
				return
			}
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Println(err)
				return
			}
			user, err := db.GetUser(db.User{
				Id: session.UserId,
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			createChatConnection("lobby", &user, conn)
		},
	}

	preHandler = func(r *web.Response) {
		r.Helper["BaseUrl"] = getBaseUrl(r.Request)
	}

	getBaseUrl = func(r web.Request) string {
		baseUrl := r.GetHeader("AppPath")
		if baseUrl == "" {
			baseUrl = "/"
		}
		return baseUrl
	}

	redirectOrRender = func(r *web.Response, url string) {
		session, _ := db.GetSession(r.Session.CookieId)
		if session.IsLoggedIn() {
			r.SetRedirect(getBaseUrl(r.Request) + url)
		} else {
			r.Render()
		}
	}
)

type WSMessage struct {
	Type string
	Data string
}

type WS_SendMessage struct {
	Message string
}

func main() {
	server := web.Server{
		Port: 8248,
		UseSessions: true,
		TemplatesDir: "templates",
		StaticFilesDir: "pub",
		PreHandler: preHandler,
		Routes: []web.Route{
			defaultRoute,
			signupRoute,
			signupSubmitRoute,
			loginRoute,
			loginSubmitRoute,
			logoutRoute,
			lobbyRoute,
			chatRoute,
		},
	}

	server.Run()
}

func createChatConnection(chatroomName string, user *db.User, ws *websocket.Conn) {
	subscription := chat.Subscribe(chatroomName, &chat.User{
		Id: user.Id,
		Username: user.Username,
	})
	defer subscription.Unsubscribe()
	waitToClose := make(chan bool)

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				break
			}
			var wsMessage WSMessage
			err = json.Unmarshal([]byte(msg), &wsMessage)
			if err != nil {
				fmt.Printf("err: %#v\n", err)
				continue
			}
			switch wsMessage.Type {
			case WS_CTS_SendMessage:
				var wsSendMessage WS_SendMessage
				err = json.Unmarshal([]byte(wsMessage.Data), &wsSendMessage)
				if err == nil {
					message := db.Message{
						Date: time.Now().Unix(),
						Message: html.EscapeString(wsSendMessage.Message),
						Chatroom: chatroomName,
						User: *user,
						UserId: user.Id,
					}
					fmt.Printf("Message: %#v\n", message)
					message.Save()
					subscription.SendMessage(&chat.Message{
						Id: message.Id,
						Date: message.Date,
						Message: message.Message,
						User: &chat.User{
							Id: message.User.Id,
							Username: message.User.Username,
						},
					})
				}
			}
		}
		waitToClose <- true
	}()

	var messageToSend *chat.Message
	var userToSend *chat.User
	for {
		select {
		case messageToSend = <-subscription.Messages:
			sendDataToWebSocket(ws, messageToSend, WS_STC_ReceiveMessage)
		case userToSend = <-subscription.UserEnter:
			sendDataToWebSocket(ws, userToSend, WS_STC_UserEnter)
		case userToSend = <-subscription.UserExit:
			sendDataToWebSocket(ws, userToSend, WS_STC_UserExit)
		case <-waitToClose:
			return
		}
	}
}

func sendDataToWebSocket(ws *websocket.Conn, value interface{}, messageType string) {
	data, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("err: %#v\n", err)
		return
	}
	msg := WSMessage{Type: messageType, Data: string(data)}
	ws.WriteJSON(msg)
}
