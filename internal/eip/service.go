package eip

import (
	"context"
	"errors"
)

// Service handles business logic for EIP cases
type Service struct {
	repo *Repository
}

// NewService creates a new EIP service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateCase creates a new EIP case
func (s *Service) CreateCase(ctx context.Context, req *CreateCaseRequest) (*EIPCase, error) {
	// Validate required fields
	if req.Side == "" {
		return nil, errors.New("side is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}

	// Validate side
	if req.Side != SideODFI && req.Side != SideRDFI {
		return nil, errors.New("side must be ODFI or RDFI")
	}

	// Validate type
	validTypes := map[string]bool{
		TypeReturnReview:    true,
		TypeNOCReview:       true,
		TypeCustomerDispute: true,
	}
	if !validTypes[req.Type] {
		return nil, errors.New("invalid case type")
	}

	eipCase := &EIPCase{
		Side:        req.Side,
		TraceNumber: req.TraceNumber,
		Status:      StatusOpen,
		Type:        req.Type,
		Notes:       req.Notes,
	}

	if err := s.repo.Create(ctx, eipCase); err != nil {
		return nil, err
	}

	return eipCase, nil
}

// GetCase retrieves an EIP case by ID
func (s *Service) GetCase(ctx context.Context, id string) (*EIPCase, error) {
	return s.repo.GetByID(ctx, id)
}

// ListCases retrieves EIP cases with optional filters
func (s *Service) ListCases(ctx context.Context, status, side, traceNumber string) ([]*EIPCase, error) {
	return s.repo.List(ctx, status, side, traceNumber)
}

// UpdateCaseStatus updates the status of an EIP case
func (s *Service) UpdateCaseStatus(ctx context.Context, id, status string) (*EIPCase, error) {
	// Validate status
	validStatuses := map[string]bool{
		StatusOpen:       true,
		StatusInProgress: true,
		StatusResolved:   true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

