package bill

import (
	"context"
	"fmt"

	"encore.app/bill/activity"
	"encore.app/bill/workflow"

	"encore.dev"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	envName       = encore.Meta().Environment.Name
	billsTaskQueue = envName + "-bills"
)

// Service is the bill service.
//
//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
}

// initService initializes the service.
func initService() (*Service, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %w", err)
	}

	w := worker.New(c, billsTaskQueue, worker.Options{})
	w.RegisterActivity(activity.FetchBillActivity)

	w.RegisterWorkflow(workflow.CreateBillWorkflow)
	w.RegisterActivity(activity.CreateBillActivity)

	w.RegisterWorkflow(workflow.ChargeBillWorkflow)
	w.RegisterActivity(activity.ChargeBillActivity)
	w.RegisterActivity(activity.CloseBillActivity)

	if err := w.Start(); err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %w", err)
	}

	return &Service{client: c, worker: w}, nil
}

// Shutdown gracefully shuts down the service.
func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}
