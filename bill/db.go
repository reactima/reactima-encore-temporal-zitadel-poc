package bill

import (
	"encore.dev/storage/sqldb"
)

// db connection with migrations applied by encore
var _ = sqldb.NewDatabase(
	"bill",
	sqldb.DatabaseConfig{
		Migrations: "./migrations",
	},
)
