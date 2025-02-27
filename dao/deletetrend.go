package dao

import (
	"context"
)

func (d *Dao) DeleteTrend(ctx context.Context) error {
	err := d.queries.DeleteTrend(ctx)

	return err
	
}
