package pubsub

import "time"

type Subscriber struct {
	EventId string
	Time    time.Time
	Listen  chan error
}
