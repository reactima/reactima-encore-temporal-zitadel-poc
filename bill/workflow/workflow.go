package workflow

import (
	"strconv"
	"time"

	"encore.app/bill/activity"
	"encore.app/bill/repository"

	"github.com/shopspring/decimal"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const (
	// QueryBill is the name of the query that returns the bill ID
	QueryBill = "query"

	// ChargeBillSignal is the name of the signal that charges a bill
	SignalChargeBill = "charge-bill"
)

// GetChargeWorkflowID returns the workflow ID for the charge workflow.
func GetChargeWorkflowID(billID int) string {
	return "bill-charges-workflow-" + strconv.Itoa(billID)
}

// CreateBillWorkflow and return Bill ID, run child workflow to charge the bill.
func CreateBillWorkflow(ctx workflow.Context, bill repository.Bill) (int, error) {
	err := workflow.
		ExecuteActivity(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: time.Second * 5,
			}),
			activity.CreateBillActivity,
			bill,
		).
		Get(ctx, &bill.ID)
	if err != nil {
		return 0, err
	}
	err = workflow.
		ExecuteActivity(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: time.Second * 5,
			}),
			activity.FetchBillActivity,
			bill,
		).
		Get(ctx, &bill)
	if err != nil {
		return 0, err
	}

	// we execute the charge workflow and abandon it until the bill is closed
	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID:         GetChargeWorkflowID(bill.ID),
		TaskQueue:          workflow.GetInfo(ctx).TaskQueueName,
		WorkflowRunTimeout: time.Second*time.Duration(bill.TimePeriod) + time.Minute,
		ParentClosePolicy:  enumspb.PARENT_CLOSE_POLICY_ABANDON,
	})

	ch := workflow.ExecuteChildWorkflow(ctx, ChargeBillWorkflow, bill).GetChildWorkflowExecution()
	if err := ch.Get(ctx, nil); err != nil {
		return 0, err
	}

	return bill.ID, nil
}

// ChargeBillWorkflow runs asynchronously and recieves signals to charge the bill.
// It will close the bill after the time period is over.
// This workflow can be queried to get the current state of the bill even when the workflow is closed.
func ChargeBillWorkflow(ctx workflow.Context, bill repository.Bill) error {
	err := workflow.SetQueryHandler(ctx, QueryBill, func() (repository.Bill, error) {
		return bill, nil
	})
	if err != nil {
		return err
	}

	selector := workflow.NewSelector(ctx)

	var closeBill bool
	closeBillTimer := workflow.NewTimer(ctx, time.Second*time.Duration(bill.TimePeriod))
	selector.AddFuture(closeBillTimer, func(future workflow.Future) {
		closeBill = true
	})

	var amount decimal.Decimal
	channel := workflow.GetSignalChannel(ctx, SignalChargeBill)
	selector.AddReceive(channel, func(c workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &amount)
		bill.Total = bill.Total.Add(amount)
	})

	for {
		selector.Select(ctx)
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
		})

		// we expect close bill by timer or charge bill
		if closeBill {
			err = workflow.
				ExecuteActivity(ctx, activity.CloseBillActivity, bill).
				Get(ctx, nil)
			if err != nil {
				workflow.GetLogger(ctx).Error("failed to close bill", "err", err)
			}
		} else {
			err = workflow.
				ExecuteActivity(ctx, activity.ChargeBillActivity, bill, amount).
				Get(ctx, nil)
			if err != nil {
				workflow.GetLogger(ctx).Error("failed to charge bill", "amount", amount, "err", err)
			}
		}

		// in any case updale the local state of the bill from the database
		// to make sure we have the latest state and it safe to query after the workflow is closed
		err = workflow.
			ExecuteActivity(ctx, activity.FetchBillActivity, bill).
			Get(ctx, &bill)
		if err != nil {
			workflow.GetLogger(ctx).Error("failed to fetch bill", "err", err)
		}

		if closeBill {
			return nil
		}
	}
}
