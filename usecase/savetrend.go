package usecase

import (
	"context"
	"log"
)

func (u *Usecase) SaveTrend(ctx context.Context) error {
	tx, err := u.dao.Begin()
	if err != nil {
		return err
	}

	err = u.dao.DeleteTrend(ctx, tx)
	if err != nil {

		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("ロールバック中にエラーが発生しました: %v", rbErr)
		}
		return err
	}

	queries, err := u.api.GetTrend()
	if err != nil {
		return err
	}

	items, err := u.api.CallPerplexityAPI(queries)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := u.dao.SaveTrend(ctx, tx, item.Name, item.Location, item.Rank, item.EndTimestamp, item.IncreasePercentage)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("ロールバック中にエラーが発生しました: %v", rbErr)
			}
			return err
		}
	}

	return err
}
