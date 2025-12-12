package odfi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for ODFI entries
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new ODFI repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const schema = `
CREATE TABLE IF NOT EXISTS odfi_entries (
	id UUID PRIMARY KEY,
	trace_number TEXT NOT NULL,
	company_name TEXT,
	sec_code TEXT,
	amount_cents BIGINT,
	status TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_odfi_entries_trace_number ON odfi_entries(trace_number);
CREATE INDEX IF NOT EXISTS idx_odfi_entries_status ON odfi_entries(status);
`

// GetSchema returns the SQL schema for ODFI tables
func GetSchema() string {
	return schema
}

// Create creates a new ODFI entry
func (r *Repository) Create(ctx context.Context, entry *ODFIEntry) error {
	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	query := `
		INSERT INTO odfi_entries (id, trace_number, company_name, sec_code, amount_cents, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.TraceNumber, entry.CompanyName, entry.SecCode,
		entry.AmountCents, entry.Status, entry.CreatedAt, entry.UpdatedAt)

	return err
}

// GetByID retrieves an ODFI entry by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*ODFIEntry, error) {
	query := `
		SELECT id, trace_number, company_name, sec_code, amount_cents, status, created_at, updated_at
		FROM odfi_entries
		WHERE id = $1
	`

	entry := &ODFIEntry{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID, &entry.TraceNumber, &entry.CompanyName, &entry.SecCode,
		&entry.AmountCents, &entry.Status, &entry.CreatedAt, &entry.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// List retrieves ODFI entries with optional filters
func (r *Repository) List(ctx context.Context, status, traceNumber string) ([]*ODFIEntry, error) {
	query := `
		SELECT id, trace_number, company_name, sec_code, amount_cents, status, created_at, updated_at
		FROM odfi_entries
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
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

	var entries []*ODFIEntry
	for rows.Next() {
		entry := &ODFIEntry{}
		err := rows.Scan(
			&entry.ID, &entry.TraceNumber, &entry.CompanyName, &entry.SecCode,
			&entry.AmountCents, &entry.Status, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// UpdateStatus updates the status of an ODFI entry
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) (*ODFIEntry, error) {
	query := `
		UPDATE odfi_entries
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, trace_number, company_name, sec_code, amount_cents, status, created_at, updated_at
	`

	entry := &ODFIEntry{}
	err := r.db.QueryRowContext(ctx, query, status, time.Now(), id).Scan(
		&entry.ID, &entry.TraceNumber, &entry.CompanyName, &entry.SecCode,
		&entry.AmountCents, &entry.Status, &entry.CreatedAt, &entry.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entry, nil
}

