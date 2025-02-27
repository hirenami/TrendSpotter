package dao

import (
	"context"
	"database/sql"
)

func (d *Dao) DeleteTrend(ctx context.Context, tx *sql.Tx) error {
	txQueries := d.WithTx(tx)

	return txQueries.DeleteTrend(ctx)

}
