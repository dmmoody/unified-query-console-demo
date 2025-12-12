package eip

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for EIP cases
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new EIP repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const schema = `
CREATE TABLE IF NOT EXISTS eip_cases (
	id UUID PRIMARY KEY,
	side TEXT NOT NULL,
	trace_number TEXT,
	status TEXT NOT NULL,
	type TEXT NOT NULL,
	notes TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_eip_cases_side ON eip_cases(side);
CREATE INDEX IF NOT EXISTS idx_eip_cases_status ON eip_cases(status);
CREATE INDEX IF NOT EXISTS idx_eip_cases_trace_number ON eip_cases(trace_number);
CREATE INDEX IF NOT EXISTS idx_eip_cases_type ON eip_cases(type);
`

// GetSchema returns the SQL schema for EIP tables
func GetSchema() string {
	return schema
}

// Create creates a new EIP case
func (r *Repository) Create(ctx context.Context, eipCase *EIPCase) error {
	eipCase.ID = uuid.New().String()
	eipCase.CreatedAt = time.Now()
	eipCase.UpdatedAt = time.Now()

	query := `
		INSERT INTO eip_cases (id, side, trace_number, status, type, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		eipCase.ID, eipCase.Side, eipCase.TraceNumber,
		eipCase.Status, eipCase.Type, eipCase.Notes,
		eipCase.CreatedAt, eipCase.UpdatedAt)

	return err
}

// GetByID retrieves an EIP case by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*EIPCase, error) {
	query := `
		SELECT id, side, trace_number, status, type, notes, created_at, updated_at
		FROM eip_cases
		WHERE id = $1
	`

	eipCase := &EIPCase{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&eipCase.ID, &eipCase.Side, &eipCase.TraceNumber,
		&eipCase.Status, &eipCase.Type, &eipCase.Notes,
		&eipCase.CreatedAt, &eipCase.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return eipCase, nil
}

// List retrieves EIP cases with optional filters
func (r *Repository) List(ctx context.Context, status, side, traceNumber string) ([]*EIPCase, error) {
	query := `
		SELECT id, side, trace_number, status, type, notes, created_at, updated_at
		FROM eip_cases
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
		argNum++
	}

	if side != "" {
		query += fmt.Sprintf(" AND side = $%d", argNum)
		args = append(args, side)
		argNum++
	}

	if traceNumber != "" {
		query += fmt.Sprintf(" AND trace_number = $%d", argNum)
		args = append(args, traceNumber)
		argNum++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []*EIPCase
	for rows.Next() {
		eipCase := &EIPCase{}
		err := rows.Scan(
			&eipCase.ID, &eipCase.Side, &eipCase.TraceNumber,
			&eipCase.Status, &eipCase.Type, &eipCase.Notes,
			&eipCase.CreatedAt, &eipCase.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cases = append(cases, eipCase)
	}

	return cases, rows.Err()
}

// UpdateStatus updates the status of an EIP case
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) (*EIPCase, error) {
	query := `
		UPDATE eip_cases
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, side, trace_number, status, type, notes, created_at, updated_at
	`

	eipCase := &EIPCase{}
	err := r.db.QueryRowContext(ctx, query, status, time.Now(), id).Scan(
		&eipCase.ID, &eipCase.Side, &eipCase.TraceNumber,
		&eipCase.Status, &eipCase.Type, &eipCase.Notes,
		&eipCase.CreatedAt, &eipCase.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return eipCase, nil
}

