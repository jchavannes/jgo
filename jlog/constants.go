package jlog

type LogLevel string

const (
	DEBUG   LogLevel = "debug"
	DEFAULT LogLevel = "default"
	TRACE   LogLevel = "trace"
	ERROR   LogLevel = "error"
)

// Order of log level. If log level is >= the current log level it will be outputted. e.g.
// - If current log level set to trace, a message marked as default WILL be outputted
// - If current log level set to default, a message marked as trace WILL NOT be outputted
var LevelOrder []LogLevel = []LogLevel{
	DEFAULT,
	DEBUG,
	TRACE,
}

func ShouldLog(messageLogLevel LogLevel, writerLogLevel LogLevel) bool {
	var messageLogLevelFound bool
	var writerLogLevelFound bool
	for _, orderLogLevel := range LevelOrder {
		if orderLogLevel == messageLogLevel {
			messageLogLevelFound = true
		}
		if orderLogLevel == writerLogLevel {
			writerLogLevelFound = true
		}
		// If we've already reached the maximum writer log level but not the message level, then don't log
		// e.g. writerLogLevel = DEFAULT, messageLogLevel = DEBUG (don't log)
		if writerLogLevelFound == true && messageLogLevelFound == false {
			return false
		}
	}
	return true
}
