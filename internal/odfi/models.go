package odfi

import (
	"time"
)

// ODFIEntry represents an origination ACH entry
type ODFIEntry struct {
	ID          string    `json:"id"`
	TraceNumber string    `json:"trace_number"`
	CompanyName string    `json:"company_name"`
	SecCode     string    `json:"sec_code"`
	AmountCents int64     `json:"amount_cents"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateEntryRequest represents the request to create an ODFI entry
type CreateEntryRequest struct {
	TraceNumber string `json:"trace_number"`
	CompanyName string `json:"company_name"`
	SecCode     string `json:"sec_code"`
	AmountCents int64  `json:"amount_cents"`
}

// UpdateStatusRequest represents the request to update an entry status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// Status constants
const (
	StatusPending   = "PENDING"
	StatusSent      = "SENT"
	StatusCancelled = "CANCELLED"
)

