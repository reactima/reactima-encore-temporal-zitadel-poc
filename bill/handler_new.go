package bill

import (
	"context"
	"fmt"

	"encore.app/bill/repository"
	"encore.app/bill/workflow"

	"encore.dev/rlog"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type NewRequest struct {
	Customer            int  `json:"customer"`
	TimePeriodInSeconds uint `json:"time_period_in_seconds"`
}

type NewResponse struct {
	BillID int `json:"bill_id"`
}

//encore:api public method=POST path=/bill/new
func (s *Service) New(ctx context.Context, request *NewRequest) (*NewResponse, error) {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("bill-workflow-%d-%s", request.Customer, uuid.New().String()),
		TaskQueue: billsTaskQueue,
	}

	bill := repository.Bill{
		CustomerID: request.Customer,
		TimePeriod: request.TimePeriodInSeconds,
	}
	we, err := s.client.ExecuteWorkflow(
		context.Background(),
		options,
		workflow.CreateBillWorkflow,
		bill,
	)
	if err != nil {
		return nil, err
	}
	rlog.Error("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	var response NewResponse
	if err := we.Get(ctx, &response.BillID); err != nil {
		return nil, err
	}

	return &response, nil
}
