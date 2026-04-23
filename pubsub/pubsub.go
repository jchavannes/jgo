package pubsub

import (
	"sync"
	"time"

	"github.com/jchavannes/jgo/jerr"
)

type PubSub struct {
	mu          sync.Mutex
	subscribers []*Subscriber
	events      []*Event

	SubscriberTimeout time.Duration
	EventTimeout      time.Duration
}

func (p *PubSub) Subscribe(eventId string) *Subscriber {
	p.mu.Lock()
	defer p.mu.Unlock()
	var subscriber = &Subscriber{
		EventId: eventId,
		Time:    time.Now(),
		Listen:  make(chan error, 1),
	}
	for _, event := range p.events {
		if event.Id == eventId {
			subscriber.Listen <- nil
			return subscriber
		}
	}
	p.subscribers = append(p.subscribers, subscriber)
	return subscriber
}

func (p *PubSub) Unsubscribe(subscriber *Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, activeSubscriber := range p.subscribers {
		if activeSubscriber == subscriber {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			return
		}
	}
}

func (p *PubSub) Publish(eventId string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, event := range p.events {
		if event.Id == eventId {
			return
		}
	}
	p.events = append(p.events, &Event{
		Id:   eventId,
		Time: time.Now(),
	})
	for i := 0; i < len(p.subscribers); i++ {
		if p.subscribers[i].EventId == eventId {
			p.subscribers[i].Listen <- nil
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
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
	p.mu.Lock()
	defer p.mu.Unlock()
	eventTimeout := time.Now().Add(-p.EventTimeout)
	for i := 0; i < len(p.events); i++ {
		if p.events[i].Time.Before(eventTimeout) {
			p.events = append(p.events[:i], p.events[i+1:]...)
			i--
		}
	}
	if p.SubscriberTimeout > 0 {
		timeout := time.Now().Add(-p.SubscriberTimeout)
		for i := 0; i < len(p.subscribers); i++ {
			if p.subscribers[i].Time.Before(timeout) {
				p.subscribers[i].Listen <- jerr.New("error pub sub timeout reached")
				p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
				i--
			}
		}
	}
}
