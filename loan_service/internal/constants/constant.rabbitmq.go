package constants

const (
	ExchangeTypeDirect = "direct"

	EmailExchange = "email_exchange"
	LogExchange   = "log_exchange"

	OTPQueue                = "otp_code"
	LoanNotificationQueue   = "loan_notification"
	ReturnNotificationQueue = "return_notification"
	LogQueue                = "log_queue"

	LogServiceAuth = "auth-service"

	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"
)
