package models

import "time"

type LoanNotificationMessage struct {
	RequestID string    `json:"X-Correlation-ID"` // for logging purpose
	Email     string    `json:"email"`
	Book      string    `json:"book"` // book title
	Due       time.Time `json:"due"`
}

type ReturnNotificationMessage struct {
	RequestID string `json:"X-Correlation-ID"` // for logging purpose
	Email     string `json:"email"`
	Book      string `json:"book"` // book title
}
