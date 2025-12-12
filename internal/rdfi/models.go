package rdfi

import (
	"time"
)

// RDFIEntry represents a receiving ACH entry
type RDFIEntry struct {
	ID           string    `json:"id"`
	TraceNumber  string    `json:"trace_number"`
	ReceiverName string    `json:"receiver_name"`
	AmountCents  int64     `json:"amount_cents"`
	Status       string    `json:"status"`
	ReturnReason string    `json:"return_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateEntryRequest represents the request to create an RDFI entry
type CreateEntryRequest struct {
	TraceNumber  string `json:"trace_number"`
	ReceiverName string `json:"receiver_name"`
	AmountCents  int64  `json:"amount_cents"`
}

// ReturnRequest represents the request to return an entry
type ReturnRequest struct {
	Reason string `json:"reason"`
}

// Status constants
const (
	StatusReceived = "RECEIVED"
	StatusPosted   = "POSTED"
	StatusReturned = "RETURNED"
)

