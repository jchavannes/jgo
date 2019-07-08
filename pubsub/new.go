package pubsub

import "time"

const (
	defaultExpireTime = 5 * time.Second
	defaultTimeout    = 30 * time.Second
)

func New() *PubSub {
	var ps = PubSub{
		ExpireTime: defaultExpireTime,
		Timeout:    defaultTimeout,
	}
	ps.initExpireChecks()
	return &ps
}
