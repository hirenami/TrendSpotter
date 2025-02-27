package dao

import (
	"context"

	"github.com/hirenami/TrendSpotter/sqlc"
)

func (d *Dao) SaveTrend(ctx context.Context, trendName, trendLocation string, trendRank int32) error{

	args := sqlc.SaveTrendParams{
		TrendsName:     trendName,
		TrendsLocation: trendLocation,
		TrendsRank:     trendRank,
	}
	err := d.queries.SaveTrend(ctx, args)

	return err
}
