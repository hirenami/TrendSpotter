package dao

import (
	"context"

	"github.com/hirenami/TrendSpotter/sqlc"
)

func (d *Dao) SaveTrend(ctx context.Context, trendName, trendLocation string, trendRank, trendEndtimestamp, trendIncreasepercentage int32) error {

	args := sqlc.SaveTrendParams{
		TrendsName:               trendName,
		TrendsLocation:           trendLocation,
		TrendsRank:               trendRank,
		TrendsEndtimestamp:       trendEndtimestamp,
		TrendsIncreasePercentage: trendIncreasepercentage,
	}
	err := d.queries.SaveTrend(ctx, args)

	return err
}
