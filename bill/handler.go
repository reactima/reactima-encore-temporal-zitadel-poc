package bill

import (
	"context"

	"encore.app/bill/repository"
	"encore.app/bill/workflow"
)

//encore:api public method=GET path=/bill/:billID
func (s *Service) Get(ctx context.Context, billID int) (*repository.Bill, error) {
	bill := repository.Bill{ID: billID}

	wID := workflow.GetChargeWorkflowID(billID)

	query, err := s.client.QueryWorkflow(ctx, wID, "", workflow.QueryBill)
	if err == nil {
		if err := query.Get(&bill); err != nil {
			return nil, err
		}
	} else {
		if err := bill.Fetch(ctx); err != nil {
			return nil, err
		}
	}

	return &bill, nil
}
