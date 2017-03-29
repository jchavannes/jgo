package db

import (
	"fmt"
)

type Message struct {
	Id       uint `gorm:"primary_key"`
	Date     int64
	Chatroom string
	Message  string
	User     User
	UserId   uint
}

func (m *Message) Save() {
	save(m)
}

func GetRecentMessages(chatroom string, count int) []*Message {
	var messages []*Message
	db, err := getDb()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return messages
	}
	db.Order("id desc").Limit(count).Find(&messages, &Message{Chatroom: chatroom})
	for _, message := range messages {
		if message.UserId > 0 {
			message.User, _ = GetUser(User{
				Id: message.UserId,
			})
		}
	}
	// Ordered by id desc to get most recent messages, but then this function needs to return ordered asc
	for i, j := 0, len(messages) - 1; i < j; i, j = i + 1, j - 1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages
}
