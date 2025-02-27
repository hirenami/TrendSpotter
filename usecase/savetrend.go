package usecase

import (
	"context"
)

func (u *Usecase) SaveTrend(ctx context.Context) error {
	queries, err := u.api.GetTrend()
	if err != nil {
		return err
	}

	items, err := u.api.CallPerplexityAPI(queries)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := u.dao.SaveTrend(ctx, item.Name, item.Location, item.Rank, item.EndTimestamp, item.IncreasePercentage)
		if err != nil {
			return err
		}
	}

	return err
}
