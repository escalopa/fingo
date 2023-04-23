package mygrpc

import (
	"context"

	"github.com/escalopa/fingo/pb"
	oteltracer "github.com/escalopa/fingo/token/internal/adapters/tracer"
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
	ctx, span := oteltracer.Tracer().Start(ctx, "ValidateToken")
	defer span.End()
	id, err := h.uc.TokenValidate.Execute(ctx, application.TokenValidateParams{AccessToken: req.AccessToken})
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{UserId: id.String()}, nil
}
