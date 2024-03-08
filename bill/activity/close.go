package activity

import (
	"context"

	"encore.app/bill/repository"
)

// CloseBillActivity close a bill for the charges.
func CloseBillActivity(ctx context.Context, bill repository.Bill) error {
	return bill.Close(ctx)
}
