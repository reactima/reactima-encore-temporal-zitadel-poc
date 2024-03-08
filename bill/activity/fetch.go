package activity

import (
	"context"

	"encore.app/bill/repository"
)

// FetchBillActivity fetches a bill from the database and returns the Bill ID.
func FetchBillActivity(ctx context.Context, bill repository.Bill) (repository.Bill, error) {
	if err := (&bill).Fetch(ctx); err != nil {
		return bill, err
	}
	return bill, nil
}
