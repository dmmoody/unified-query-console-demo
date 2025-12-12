package ledger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for ledger entries
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new ledger repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const schema = `
CREATE TABLE IF NOT EXISTS ledger_entries (
	id UUID PRIMARY KEY,
	ach_side TEXT NOT NULL,
	trace_number TEXT,
	amount_cents BIGINT NOT NULL,
	direction TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_ach_side ON ledger_entries(ach_side);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_trace_number ON ledger_entries(trace_number);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_direction ON ledger_entries(direction);
`

// GetSchema returns the SQL schema for ledger tables
func GetSchema() string {
	return schema
}

// Create creates a new ledger entry
func (r *Repository) Create(ctx context.Context, entry *LedgerEntry) error {
	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()

	query := `
		INSERT INTO ledger_entries (id, ach_side, trace_number, amount_cents, direction, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.AchSide, entry.TraceNumber,
		entry.AmountCents, entry.Direction, entry.Description,
		entry.CreatedAt)

	return err
}

// List retrieves ledger entries with optional filters
func (r *Repository) List(ctx context.Context, achSide, traceNumber string) ([]*LedgerEntry, error) {
	query := `
		SELECT id, ach_side, trace_number, amount_cents, direction, description, created_at
		FROM ledger_entries
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if achSide != "" {
		query += fmt.Sprintf(" AND ach_side = $%d", argNum)
		args = append(args, achSide)
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

	var entries []*LedgerEntry
	for rows.Next() {
		entry := &LedgerEntry{}
		err := rows.Scan(
			&entry.ID, &entry.AchSide, &entry.TraceNumber,
			&entry.AmountCents, &entry.Direction, &entry.Description,
			&entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// GetBalances calculates total debits, credits, and net balance
func (r *Repository) GetBalances(ctx context.Context) (*BalanceResponse, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN direction = 'DEBIT' THEN amount_cents ELSE 0 END), 0) as total_debits,
			COALESCE(SUM(CASE WHEN direction = 'CREDIT' THEN amount_cents ELSE 0 END), 0) as total_credits
		FROM ledger_entries
	`

	var totalDebits, totalCredits int64
	err := r.db.QueryRowContext(ctx, query).Scan(&totalDebits, &totalCredits)
	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
		NetBalance:   totalCredits - totalDebits,
	}, nil
}

