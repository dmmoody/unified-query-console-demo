package odfi

import (
	"context"
	"errors"
)

// Service handles business logic for ODFI entries
type Service struct {
	repo *Repository
}

// NewService creates a new ODFI service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateEntry creates a new ODFI entry
func (s *Service) CreateEntry(ctx context.Context, req *CreateEntryRequest) (*ODFIEntry, error) {
	if req.TraceNumber == "" {
		return nil, errors.New("trace_number is required")
	}

	entry := &ODFIEntry{
		TraceNumber: req.TraceNumber,
		CompanyName: req.CompanyName,
		SecCode:     req.SecCode,
		AmountCents: req.AmountCents,
		Status:      StatusPending,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetEntry retrieves an ODFI entry by ID
func (s *Service) GetEntry(ctx context.Context, id string) (*ODFIEntry, error) {
	return s.repo.GetByID(ctx, id)
}

// ListEntries retrieves ODFI entries with optional filters
func (s *Service) ListEntries(ctx context.Context, status, traceNumber string) ([]*ODFIEntry, error) {
	return s.repo.List(ctx, status, traceNumber)
}

// UpdateEntryStatus updates the status of an ODFI entry
func (s *Service) UpdateEntryStatus(ctx context.Context, id string, status string) (*ODFIEntry, error) {
	// Validate status
	validStatuses := map[string]bool{
		StatusPending:   true,
		StatusSent:      true,
		StatusCancelled: true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

