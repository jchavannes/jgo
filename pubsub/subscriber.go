package pubsub

type Subscriber struct {
	EventId string
	Listen  chan bool
}
