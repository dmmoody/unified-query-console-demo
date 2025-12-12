package ledger

import (
	"context"
	"errors"
)

// Service handles business logic for ledger entries
type Service struct {
	repo *Repository
}

// NewService creates a new ledger service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreatePosting creates a new ledger posting
func (s *Service) CreatePosting(ctx context.Context, req *CreatePostingRequest) (*LedgerEntry, error) {
	// Validate required fields
	if req.AchSide == "" {
		return nil, errors.New("ach_side is required")
	}
	if req.Direction == "" {
		return nil, errors.New("direction is required")
	}

	// Validate direction
	if req.Direction != DirectionDebit && req.Direction != DirectionCredit {
		return nil, errors.New("direction must be DEBIT or CREDIT")
	}

	// Validate ach_side
	if req.AchSide != SideODFI && req.AchSide != SideRDFI {
		return nil, errors.New("ach_side must be ODFI or RDFI")
	}

	entry := &LedgerEntry{
		AchSide:     req.AchSide,
		TraceNumber: req.TraceNumber,
		AmountCents: req.AmountCents,
		Direction:   req.Direction,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// ListPostings retrieves ledger postings with optional filters
func (s *Service) ListPostings(ctx context.Context, achSide, traceNumber string) ([]*LedgerEntry, error) {
	return s.repo.List(ctx, achSide, traceNumber)
}

// GetBalances calculates and returns balance information
func (s *Service) GetBalances(ctx context.Context) (*BalanceResponse, error) {
	return s.repo.GetBalances(ctx)
}

