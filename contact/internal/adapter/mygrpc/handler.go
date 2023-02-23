package mygrpc

import (
	"context"

	"github.com/escalopa/gofly/contact/internal/application"
	"github.com/escalopa/gofly/pb"
)

type Handler struct {
	uc *application.UseCases
	pb.EmailServiceServer
}

func New(uc *application.UseCases) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) SendCode(ctx context.Context, req *pb.SendCodeRequest) (*pb.SendCodeResponse, error) {
	err := h.uc.SendCode.Execute(ctx, application.SendCodeCommandParam{Email: req.Email})
	if err != nil {
		return nil, err
	}
	return &pb.SendCodeResponse{
		Response: &pb.BasicResponse{
			Status:  200,
			Message: "Successfully sent code",
		},
	}, nil
}

func (h *Handler) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.VerifyCodeResponse, error) {
	err := h.uc.VerifyCode.Execute(ctx, application.VerifyCodeCommandParam{Email: req.Email, Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyCodeResponse{
		Response: &pb.BasicResponse{
			Status:  200,
			Message: "Successfully verified code",
		},
	}, nil
}
