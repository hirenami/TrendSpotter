package dao

import (
	"database/sql"

	"github.com/hirenami/TrendSpotter/sqlc"
)

type Dao struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewDao コンストラクタ
func NewDao(db *sql.DB, queries *sqlc.Queries) *Dao {
	return &Dao{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (d *Dao) WithTx(tx *sql.Tx) *sqlc.Queries {
	return d.queries.WithTx(tx)
}

// Begin メソッドの実装
func (d *Dao) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}
