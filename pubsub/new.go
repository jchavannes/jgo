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
	// Listen is buffered (size 1) so PubSub can send notifications while holding its mutex
	// without blocking. Each subscriber receives at most one notification before removal.
	Listen chan error
}

func New() *PubSub {
	var ps = PubSub{
		EventTimeout:      defaultEventTimeout,
		SubscriberTimeout: defaultSubscriberTimeout,
	}
	ps.initExpireChecks()
	return &ps
}
