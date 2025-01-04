package constants

const (
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"

	LogFieldTimeStamp      = "timestamp"
	LogFieldService        = "service"
	LogFieldXCorrelationID = "X-Correlation-ID"
	LogFieldCaller         = "caller"
	LogFieldExtra          = "extra"

	LogExchange = "log_exchange"
	LogQueue    = "log_queue"
)
