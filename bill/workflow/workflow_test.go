package workflow

import (
	"testing"

	"encore.app/bill/activity"
	"encore.app/bill/repository"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func Test_CreateBillWorkflowSuccess(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(ChargeBillWorkflow)

	env.OnActivity(
		activity.CreateBillActivity,
		mock.Anything,
		repository.Bill{CustomerID: 1, Total: decimal.NewFromInt(0)},
	).Return(1, nil)

	env.OnActivity(
		activity.ChargeBillActivity,
		mock.Anything,
		repository.Bill{
			ID: 1,
			CustomerID: 1,
		}, decimal.NewFromFloat(0.22),
	)

	env.OnActivity(
		activity.FetchBillActivity,
		mock.Anything,
		repository.Bill{
			ID: 1,
			CustomerID: 1,
			Total: decimal.NewFromInt(0),
		},
	).Return(
		repository.Bill{
			ID: 1,
			CustomerID: 1,
			Total: decimal.NewFromFloat(0.22),
			Charges: []repository.BillCharge{
				{ID: 1, Amount: decimal.NewFromFloat(0.22)},
			},
		},
		nil,
	)


	env.ExecuteWorkflow(CreateBillWorkflow, repository.Bill{CustomerID: 1})
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
