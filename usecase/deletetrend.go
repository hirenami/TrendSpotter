package usecase

import "context"

func (u *Usecase) DeleteTrend(ctx context.Context)error {
	err :=u.dao.DeleteTrend(ctx)

	return err
}
