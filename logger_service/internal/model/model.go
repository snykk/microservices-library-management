package model

import "time"

// LogMessage is the structure for the log message
type LogMessage struct {
	Timestamp      time.Time              `json:"timestamp" bson:"timestamp"`
	Service        string                 `json:"service" bson:"service"`
	Level          string                 `json:"level" bson:"level"`
	XCorrelationID string                 `json:"X-Correlation-ID" bson:"X-Correlation-ID"`
	Caller         string                 `json:"caller" bson:"caller"`
	Message        string                 `json:"message" bson:"message"`
	Error          *string                `json:"error,omitempty" bson:"error,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty" bson:"extra,omitempty"`
}
