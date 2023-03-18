package mygrpc

import (
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/user/internal/application"
)

type Handler struct {
	uc *application.UseCases
	pb.UnimplementedUserServiceServer
}

func New(uc *application.UseCases) *Handler {
	return &Handler{uc: uc}
}
