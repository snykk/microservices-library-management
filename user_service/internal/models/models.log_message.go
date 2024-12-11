package models

import "time"

// LogMessage is the structure for the log message
type LogMessage struct {
	Timestamp      time.Time              `json:"timestamp"`
	Service        string                 `json:"service"`
	Level          string                 `json:"level"`
	XCorrelationID string                 `json:"X-Correlation-ID"`
	Caller         string                 `json:"caller"`
	Message        string                 `json:"message"`
	Error          string                 `json:"error,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}
