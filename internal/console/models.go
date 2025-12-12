package console

// UnifiedAchItem represents a unified view of ACH entries from ODFI/RDFI
type UnifiedAchItem struct {
	Side        string `json:"side"`          // "ODFI" or "RDFI"
	Source      string `json:"source"`        // "odfi", "rdfi"
	EntryID     string `json:"entry_id"`
	TraceNumber string `json:"trace_number"`
	AmountCents int64  `json:"amount_cents"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`    // For sorting
	Extra       any    `json:"extra,omitempty"` // Optional service-specific fields
}

// ODFIEntry represents an ODFI entry from the ODFI service
type ODFIEntry struct {
	ID          string `json:"id"`
	TraceNumber string `json:"trace_number"`
	CompanyName string `json:"company_name"`
	SecCode     string `json:"sec_code"`
	AmountCents int64  `json:"amount_cents"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateODFIEntryRequest represents request to create ODFI entry
type CreateODFIEntryRequest struct {
	TraceNumber string `json:"trace_number"`
	CompanyName string `json:"company_name"`
	SecCode     string `json:"sec_code"`
	AmountCents int64  `json:"amount_cents"`
}

// UpdateODFIStatusRequest represents request to update ODFI status
type UpdateODFIStatusRequest struct {
	Status string `json:"status"`
}

// RDFIEntry represents an RDFI entry from the RDFI service
type RDFIEntry struct {
	ID           string `json:"id"`
	TraceNumber  string `json:"trace_number"`
	ReceiverName string `json:"receiver_name"`
	AmountCents  int64  `json:"amount_cents"`
	Status       string `json:"status"`
	ReturnReason string `json:"return_reason,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateRDFIEntryRequest represents request to create RDFI entry
type CreateRDFIEntryRequest struct {
	TraceNumber  string `json:"trace_number"`
	ReceiverName string `json:"receiver_name"`
	AmountCents  int64  `json:"amount_cents"`
}

// ReturnRequest represents a request to return an entry
type ReturnRequest struct {
	Reason string `json:"reason"`
}

// LedgerEntry represents a ledger posting
type LedgerEntry struct {
	ID          string `json:"id"`
	AchSide     string `json:"ach_side"`
	TraceNumber string `json:"trace_number"`
	AmountCents int64  `json:"amount_cents"`
	Direction   string `json:"direction"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// CreateLedgerPostingRequest represents request to create ledger posting
type CreateLedgerPostingRequest struct {
	AchSide     string `json:"ach_side"`
	TraceNumber string `json:"trace_number"`
	AmountCents int64  `json:"amount_cents"`
	Direction   string `json:"direction"`
	Description string `json:"description"`
}

// BalanceResponse represents balance calculation
type BalanceResponse struct {
	TotalDebits  int64 `json:"total_debits"`
	TotalCredits int64 `json:"total_credits"`
	NetBalance   int64 `json:"net_balance"`
}

// EIPCase represents an exception case
type EIPCase struct {
	ID          string `json:"id"`
	Side        string `json:"side"`
	TraceNumber string `json:"trace_number"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateEIPCaseRequest represents request to create EIP case
type CreateEIPCaseRequest struct {
	Side        string `json:"side"`
	TraceNumber string `json:"trace_number"`
	Type        string `json:"type"`
	Notes       string `json:"notes"`
}

// UpdateEIPCaseStatusRequest represents request to update case status
type UpdateEIPCaseStatusRequest struct {
	Status string `json:"status"`
}

