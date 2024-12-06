package models

import "time"

type LoanNotificationMessage struct {
	Email string    `json:"email"`
	Book  string    `json:"book"` // book title
	Due   time.Time `json:"due"`
}

type ReturnNotificationMessage struct {
	Email string `json:"email"`
	Book  string `json:"book"` // book title
}
