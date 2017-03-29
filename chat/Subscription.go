package chat

import (
	"fmt"
)

type Subscription struct {
	Id        uint
	Messages  chan *Message
	UserEnter chan *User
	UserExit  chan *User
	chatroom  *Chatroom
	user      *User
}

func (s *Subscription) Unsubscribe() {
	fmt.Printf("Unsubscribing from chatroom (Subscription.Id: %d).\n", s.Id)
	delete(s.chatroom.subscriptions, s.Id)
	for _, subscription := range s.chatroom.subscriptions {
		if subscription.user == s.user && subscription != s {
			return
		}
	}
	for _, subscription := range s.chatroom.subscriptions {
		fmt.Printf("Sending user exit (User.Id: %d) to subscriber (Subscription.Id: %d).\n", s.user.Id, subscription.Id)
		go func(subscription *Subscription) {
			subscription.UserExit <- s.user
		}(subscription)
	}
}

func (s *Subscription) SendMessage(message *Message) {
	chatroom := s.chatroom
	chatroom.recentMessages = append(chatroom.recentMessages, message)
	if len(chatroom.recentMessages) > chatroom.maxRecentMessages {
		chatroom.recentMessages = chatroom.recentMessages[len(chatroom.recentMessages) - chatroom.maxRecentMessages:]
	}
	for _, subscription := range chatroom.subscriptions {
		go func(subscription *Subscription) {
			fmt.Printf("Sending message (%s) to subscriber (Subscription.Id: %d).\n", message, subscription.Id)
			subscription.Messages <- message
		}(subscription)
	}
}

func Subscribe(chatroomName string, user *User) *Subscription {
	chatroom := getChatRoom(chatroomName)
	chatroom.incrementer++
	subscription := &Subscription{
		Id: chatroom.incrementer,
		Messages: make(chan *Message),
		UserEnter: make(chan *User),
		UserExit: make(chan *User),
		chatroom: chatroom,
		user: user,
	}
	fmt.Printf("Subscribing to chatroom (Subscrition.Id: %d).\n", subscription.Id)
	chatroom.subscriptions[subscription.Id] = subscription
	go func() {
		for _, messageToSend := range chatroom.recentMessages {
			subscription.Messages <- messageToSend
		}
	}()
	go func() {
		users := chatroom.GetUsers()
		for _, userToSend := range users {
			fmt.Printf("Sending user enter (User.Id: %d) to subscriber (Subscription.Id: %d).\n", userToSend.Id, subscription.Id)
			subscription.UserEnter <- userToSend
		}
	}()
	go func() {
		for _, subscriptionToBeNotified := range chatroom.subscriptions {
			if subscriptionToBeNotified.user != user {
				fmt.Printf("Sending user enter (User.Id: %d) to subscriber (Subscription.Id: %d).\n", user.Id, subscriptionToBeNotified.Id)
				subscriptionToBeNotified.UserEnter <- user
			}
		}
	}()
	return subscription
}
