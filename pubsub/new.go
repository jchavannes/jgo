package pubsub

import "time"

const (
	defaultEventTimeout      = 30 * time.Second
	defaultSubscriberTimeout = 3 * time.Minute
)

type Event struct {
	Id   string
	Time time.Time
}

type Subscriber struct {
	EventId string
	Time    time.Time
	Listen  chan error
}

func New() *PubSub {
	var ps = PubSub{
		EventTimeout:      defaultEventTimeout,
		SubscriberTimeout: defaultSubscriberTimeout,
	}
	ps.initExpireChecks()
	return &ps
}
