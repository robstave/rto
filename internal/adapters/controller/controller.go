package controller

import (
	"github.com/robstave/rto/internal/domain"
)

type RTOController struct {
	service domain.RTOBLL
	test2   string
}

func NewRTOController(

	test2 string,

) *RTOController {

	service := domain.NewService(
		2.5,
		"gg",
	)

	return &RTOController{service, test2}
}
