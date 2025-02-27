package usecase

import "context"

func (u *Usecase) SaveTrend(ctx context.Context, trendName, trendLocation string, trendRank int32) error{
	err := u.dao.SaveTrend(ctx, trendName, trendLocation, trendRank)

	return err
}
