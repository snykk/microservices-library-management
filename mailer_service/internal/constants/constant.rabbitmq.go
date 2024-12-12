package constants

const (
	ExchangeTypeDirect = "direct"

	EmailExchange = "email_exchange"
	LogExchange   = "log_exchange"

	OTPQueue                = "otp_code"
	LogQueue                = "log_queue"
	LoanNotificationQueue   = "loan_notification"
	ReturnNotificationQueue = "return_notification"

	LogServiceMailer = "mailer-service"

	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"
)
