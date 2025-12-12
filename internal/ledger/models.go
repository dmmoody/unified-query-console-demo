package ledger

import (
	"time"
)

// LedgerEntry represents a ledger posting
type LedgerEntry struct {
	ID          string    `json:"id"`
	AchSide     string    `json:"ach_side"`
	TraceNumber string    `json:"trace_number"`
	AmountCents int64     `json:"amount_cents"`
	Direction   string    `json:"direction"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreatePostingRequest represents the request to create a ledger posting
type CreatePostingRequest struct {
	AchSide     string `json:"ach_side"`
	TraceNumber string `json:"trace_number"`
	AmountCents int64  `json:"amount_cents"`
	Direction   string `json:"direction"`
	Description string `json:"description"`
}

// BalanceResponse represents the balance calculation
type BalanceResponse struct {
	TotalDebits  int64 `json:"total_debits"`
	TotalCredits int64 `json:"total_credits"`
	NetBalance   int64 `json:"net_balance"`
}

// Direction constants
const (
	DirectionDebit  = "DEBIT"
	DirectionCredit = "CREDIT"
)

// Side constants
const (
	SideODFI = "ODFI"
	SideRDFI = "RDFI"
)

