package models

import (
	"context"
	"database/sql"
)

var (
	CTX     context.Context
	QUERIES *Queries
	DB      *sql.DB
)
