package rdfi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for RDFI entries
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new RDFI repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const schema = `
CREATE TABLE IF NOT EXISTS rdfi_entries (
	id UUID PRIMARY KEY,
	trace_number TEXT NOT NULL,
	receiver_name TEXT,
	amount_cents BIGINT,
	status TEXT NOT NULL,
	return_reason TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rdfi_entries_trace_number ON rdfi_entries(trace_number);
CREATE INDEX IF NOT EXISTS idx_rdfi_entries_status ON rdfi_entries(status);
`

// GetSchema returns the SQL schema for RDFI tables
func GetSchema() string {
	return schema
}

// Create creates a new RDFI entry
func (r *Repository) Create(ctx context.Context, entry *RDFIEntry) error {
	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	query := `
		INSERT INTO rdfi_entries (id, trace_number, receiver_name, amount_cents, status, return_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.TraceNumber, entry.ReceiverName,
		entry.AmountCents, entry.Status, nullString(entry.ReturnReason),
		entry.CreatedAt, entry.UpdatedAt)

	return err
}

// GetByID retrieves an RDFI entry by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*RDFIEntry, error) {
	query := `
		SELECT id, trace_number, receiver_name, amount_cents, status, return_reason, created_at, updated_at
		FROM rdfi_entries
		WHERE id = $1
	`

	entry := &RDFIEntry{}
	var returnReason sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID, &entry.TraceNumber, &entry.ReceiverName,
		&entry.AmountCents, &entry.Status, &returnReason,
		&entry.CreatedAt, &entry.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if returnReason.Valid {
		entry.ReturnReason = returnReason.String
	}

	return entry, nil
}

// List retrieves RDFI entries with optional filters
func (r *Repository) List(ctx context.Context, status, traceNumber string) ([]*RDFIEntry, error) {
	query := `
		SELECT id, trace_number, receiver_name, amount_cents, status, return_reason, created_at, updated_at
		FROM rdfi_entries
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

	var entries []*RDFIEntry
	for rows.Next() {
		entry := &RDFIEntry{}
		var returnReason sql.NullString

		err := rows.Scan(
			&entry.ID, &entry.TraceNumber, &entry.ReceiverName,
			&entry.AmountCents, &entry.Status, &returnReason,
			&entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if returnReason.Valid {
			entry.ReturnReason = returnReason.String
		}

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// Return marks an entry as returned with a reason
func (r *Repository) Return(ctx context.Context, id, reason string) (*RDFIEntry, error) {
	query := `
		UPDATE rdfi_entries
		SET status = $1, return_reason = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, trace_number, receiver_name, amount_cents, status, return_reason, created_at, updated_at
	`

	entry := &RDFIEntry{}
	var returnReason sql.NullString

	err := r.db.QueryRowContext(ctx, query, StatusReturned, reason, time.Now(), id).Scan(
		&entry.ID, &entry.TraceNumber, &entry.ReceiverName,
		&entry.AmountCents, &entry.Status, &returnReason,
		&entry.CreatedAt, &entry.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if returnReason.Valid {
		entry.ReturnReason = returnReason.String
	}

	return entry, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

