package repository

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	ctx    context.Context
	cancel context.CancelFunc
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	// use temporal transaction to rollback all changes after each test
	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Second*5)
	_, err := db.Exec(s.ctx, "BEGIN")
	require.NoError(s.T(), err)
}

func (s *UnitTestSuite) TearDownTest() {
	_, err := db.Exec(s.ctx, "ROLLBACK")
	require.NoError(s.T(), err)
	s.cancel()
}

func TestBillRepositoryMethods(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	newBillID := 0
	t.Run("create", func(t *testing.T) {
		bill, err := Bill{CustomerID: 3}.Create(ctx)
		require.NoError(t, err)
		require.NotZero(t, bill.ID)
		newBillID = bill.ID
	})

	t.Run("charge", func(t *testing.T) {
		err := Bill{ID: newBillID}.Charge(ctx, decimal.NewFromFloat(22.15))
		require.NoError(t, err)

		err = Bill{ID: newBillID}.Charge(ctx, decimal.NewFromFloat(42.05))
		require.NoError(t, err)
	})

	t.Run("fetch", func(t *testing.T) {
		expectedAmount, err := decimal.NewFromString("64.20")
		require.NoError(t, err)

		bill := &Bill{ID: newBillID}
		err = bill.Fetch(ctx)
		require.NoError(t, err)
		require.Equal(t, newBillID, bill.ID)
		require.Equal(t, 3, bill.CustomerID)
		require.True(
			t, expectedAmount.Equal(bill.Total), "total mismatch: expected 64.20, got: %v", bill.Total)
		require.False(t, bill.CreatedAt.IsZero())
		require.Nil(t, bill.ClosedAt)
		require.Len(t, bill.Charges, 2)

		require.Equal(t, "22.15", bill.Charges[0].Amount.String())
		require.Equal(t, "42.05", bill.Charges[1].Amount.String())
	})

	t.Run("close", func(t *testing.T) {
		err := Bill{ID: newBillID}.Close(ctx)
		require.NoError(t, err)

		bill := &Bill{ID: newBillID}
		err = bill.Fetch(ctx)
		require.NoError(t, err)
		require.False(t, bill.ClosedAt.IsZero())
	})
}
