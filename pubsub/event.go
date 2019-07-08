package pubsub

import "time"

type Event struct {
	Id   string
	Time time.Time
}
