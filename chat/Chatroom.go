package chat

import (
	"fmt"
)

type Chatroom struct {
	maxRecentMessages int
	recentMessages    []*Message
	subscriptions     map[uint]*Subscription
	incrementer       uint
}

func (c *Chatroom) GetUsers() []*User {
	var users []*User
	for _, subscription := range c.subscriptions {
		userFound := false
		for _, user := range users {
			if subscription.user == user {
				userFound = true
				break;
			}
		}
		if userFound {
			continue;
		}
		users = append(users, subscription.user)
	}
	return users
}

var chatRooms map[string]*Chatroom

func getChatRoom(chatroomName string) *Chatroom {
	if chatRooms == nil {
		chatRooms = make(map[string]*Chatroom)
	}
	if _, ok := chatRooms[chatroomName]; !ok {
		fmt.Printf("Initializing chatroom (%s).\n", chatroomName)
		chatRooms[chatroomName] = &Chatroom{
			maxRecentMessages: 100,
			subscriptions: make(map[uint]*Subscription),
		}
	}
	return chatRooms[chatroomName]
}
