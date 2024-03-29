package mygrpc

import (
	"context"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/token/internal/application"
)

type TokenHandler struct {
	uc *application.UseCases
	pb.UnimplementedTokenServiceServer
}

func NewTokenHandler(uc *application.UseCases) *TokenHandler {
	return &TokenHandler{uc: uc}
}

func (h *TokenHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "ValidateToken")
	defer span.End()
	id, err := h.uc.TokenValidate.Execute(ctx, application.TokenValidateParams{AccessToken: req.AccessToken})
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{UserId: id.String()}, nil
}
