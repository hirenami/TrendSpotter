package handler

import "github.com/hirenami/TrendSpotter/usecase"

type Handler struct {
	Usecase *usecase.Usecase
}

func Newhandler(usecase *usecase.Usecase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}
