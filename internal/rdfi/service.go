package rdfi

import (
	"context"
	"errors"
)

// Service handles business logic for RDFI entries
type Service struct {
	repo *Repository
}

// NewService creates a new RDFI service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateEntry creates a new RDFI entry
func (s *Service) CreateEntry(ctx context.Context, req *CreateEntryRequest) (*RDFIEntry, error) {
	if req.TraceNumber == "" {
		return nil, errors.New("trace_number is required")
	}

	entry := &RDFIEntry{
		TraceNumber:  req.TraceNumber,
		ReceiverName: req.ReceiverName,
		AmountCents:  req.AmountCents,
		Status:       StatusReceived,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetEntry retrieves an RDFI entry by ID
func (s *Service) GetEntry(ctx context.Context, id string) (*RDFIEntry, error) {
	return s.repo.GetByID(ctx, id)
}

// ListEntries retrieves RDFI entries with optional filters
func (s *Service) ListEntries(ctx context.Context, status, traceNumber string) ([]*RDFIEntry, error) {
	return s.repo.List(ctx, status, traceNumber)
}

// ReturnEntry marks an entry as returned
func (s *Service) ReturnEntry(ctx context.Context, id, reason string) (*RDFIEntry, error) {
	if reason == "" {
		return nil, errors.New("return reason is required")
	}

	return s.repo.Return(ctx, id, reason)
}

