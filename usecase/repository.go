package usecase

import (
	"github.com/hirenami/TrendSpotter/api"
	"github.com/hirenami/TrendSpotter/dao"
)

// Usecase構造体
type Usecase struct {
	dao *dao.Dao
	api *api.Api
}

// NewTestUsecase コンストラクタ
func NewUsecase(dao *dao.Dao) *Usecase {
	return &Usecase{
		dao: dao,
		api: api.NewApi(),
	}
}
