package activity

import (
	"context"

	"encore.app/bill/repository"
	"github.com/shopspring/decimal"
)

// ChargeBillActivity charges a new bill in the database and returns the Bill ID.
func ChargeBillActivity(ctx context.Context, bill repository.Bill, amount decimal.Decimal) error {
	return bill.Charge(ctx, amount)
}
