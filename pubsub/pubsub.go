package pubsub

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type PubSub struct {
	Subscribers []*Subscriber
	Events      []*Event

	SubscriberTimeout time.Duration
	EventTimeout      time.Duration
}

func (p *PubSub) Subscribe(eventId string) *Subscriber {
	var subscriber = &Subscriber{
		EventId: eventId,
		Time:    time.Now(),
		Listen:  make(chan error),
	}
	p.Subscribers = append(p.Subscribers, subscriber)
	for _, event := range p.Events {
		if event.Id == eventId {
			go func() { subscriber.Listen <- nil }()
			p.Unsubscribe(subscriber)
		}
	}
	return subscriber

}

func (p *PubSub) Unsubscribe(subscriber *Subscriber) {
	for i, activeSubscriber := range p.Subscribers {
		if activeSubscriber == subscriber {
			p.Subscribers = append(p.Subscribers[:i], p.Subscribers[i+1:]...)
			return
		}
	}
}

func (p *PubSub) Publish(eventId string) {
	for _, event := range p.Events {
		if event.Id == eventId {
			return
		}
	}
	p.Events = append(p.Events, &Event{
		Id:   eventId,
		Time: time.Now(),
	})
	for i := 0; i < len(p.Subscribers); i++ {
		if p.Subscribers[i].EventId == eventId {
			go func(sub *Subscriber) { sub.Listen <- nil }(p.Subscribers[i])
			p.Subscribers = append(p.Subscribers[:i], p.Subscribers[i+1:]...)
			i--
		}
	}
}

func (p *PubSub) initExpireChecks() {
	go func() {
		for {
			time.Sleep(time.Second)
			p.checkExpired()
		}
	}()
}

func (p *PubSub) checkExpired() {
	eventTimeout := time.Now().Add(-p.EventTimeout)
	for i := 0; i < len(p.Events); i++ {
		if p.Events[i].Time.Before(eventTimeout) {
			p.Events = append(p.Events[:i], p.Events[i+1:]...)
			i--
		}
	}
	if p.SubscriberTimeout > 0 {
		timeout := time.Now().Add(-p.SubscriberTimeout)
		for i := 0; i < len(p.Subscribers); i++ {
			if p.Subscribers[i].Time.Before(timeout) {
				go func(sub *Subscriber) { sub.Listen <- jerr.New("error pub sub timeout reached") }(p.Subscribers[i])
				p.Subscribers = append(p.Subscribers[:i], p.Subscribers[i+1:]...)
				i--
			}
		}
	}
}
