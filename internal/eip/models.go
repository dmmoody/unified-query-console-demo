package eip

import (
	"time"
)

// EIPCase represents an exception/investigation case
type EIPCase struct {
	ID          string    `json:"id"`
	Side        string    `json:"side"`
	TraceNumber string    `json:"trace_number"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCaseRequest represents the request to create a case
type CreateCaseRequest struct {
	Side        string `json:"side"`
	TraceNumber string `json:"trace_number"`
	Type        string `json:"type"`
	Notes       string `json:"notes"`
}

// UpdateStatusRequest represents the request to update case status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// Status constants
const (
	StatusOpen       = "OPEN"
	StatusInProgress = "IN_PROGRESS"
	StatusResolved   = "RESOLVED"
)

// Type constants
const (
	TypeReturnReview     = "RETURN_REVIEW"
	TypeNOCReview        = "NOC_REVIEW"
	TypeCustomerDispute  = "CUSTOMER_DISPUTE"
)

// Side constants
const (
	SideODFI = "ODFI"
	SideRDFI = "RDFI"
)

