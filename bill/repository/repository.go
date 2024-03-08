package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
	"github.com/shopspring/decimal"
)

var db = sqldb.Named("bill")

// BillCharge is a database model for the bill_charges table.
type BillCharge struct {
	ID        int             `json:"id"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}

// Bill is a database model for the bill table.
type Bill struct {
	ID         int             `json:"id"`
	CustomerID int             `json:"customer"`
	TimePeriod uint            `json:"time_period"`
	Total      decimal.Decimal `json:"total"`
	CreatedAt  time.Time       `json:"created_at"`
	ClosedAt   *time.Time      `json:"closed_at"`

	Charges []BillCharge `json:"charges"`
}

// Create creates a new bill in the database.
func (b Bill) Create(ctx context.Context) (*Bill, error) {
	err := db.QueryRow(
		ctx,
		(`INSERT INTO bills (customer_id, time_period) `+
			`VALUES ($1, $2) RETURNING id, created_at;`),
		b.CustomerID, b.TimePeriod,
	).Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}

	return &b, nil
}

// Fetch fetches a bill from the database.
func (b *Bill) Fetch(ctx context.Context) error {
	var closedAt sql.NullTime
	err := db.QueryRow(
		ctx,
		(`SELECT`+
			` b.customer_id, b.time_period, b.created_at, b.closed_at,`+
			` (SELECT COALESCE(sum(amount), 0) AS total FROM bill_charges WHERE bill_id=$1)`+
			` FROM bills AS b WHERE id = $1`),
		b.ID,
	).Scan(&b.CustomerID, &b.TimePeriod, &b.CreatedAt, &closedAt, &b.Total)
	if err != nil {
		return fmt.Errorf("failed to create bill: %w", err)
	}
	if closedAt.Valid {
		b.ClosedAt = &closedAt.Time
		b.Total = b.Total.RoundUp(2) // As final bill pay amount need to be rounded up to 2 decimal places.
	}

	stmt, err := db.Query(
		ctx,
		(`SELECT` +
			` c.id, c.bill_id, c.amount, c.created_at AS charge_create_at` +
			` FROM bill_charges AS c WHERE c.bill_id = $1 ORDER BY id;`),
		b.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to create bill: %w", err)
	}
	defer stmt.Close()

	b.Charges = make([]BillCharge, 0)
	for stmt.Next() {
		var c BillCharge
		if err := stmt.Scan(&c.ID, nil, &c.Amount, &c.CreatedAt); err != nil {
			return fmt.Errorf("fetching bill charges: %w", err)
		}
		b.Charges = append(b.Charges, c)
	}

	return nil
}

// Close closes a bill set closed date to now.
func (b Bill) Close(ctx context.Context) error {
	if b.ID == 0 {
		return fmt.Errorf("bill ID can't be 0 when close")
	}
	_, err := db.Exec(
		ctx,
		`UPDATE bills SET closed_at = NOW() WHERE id = $1;`,
		b.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to close bill: %w", err)
	}
	return nil
}

// Charge adds a charge to a bill.
func (b Bill) Charge(ctx context.Context, amount decimal.Decimal) error {
	if b.ID == 0 {
		return fmt.Errorf("bill ID can't be 0 when charge")
	}
	_, err := db.Exec(
		ctx,
		`INSERT INTO bill_charges (bill_id, amount) VALUES ($1, $2);`,
		b.ID, amount,
	)
	return err
}
