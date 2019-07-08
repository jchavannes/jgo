package pubsub

var pubSubInitialized bool

func Init() {
	if pubSubInitialized {
		return
	}
	pubSubInitialized = true
}
