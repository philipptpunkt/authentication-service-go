package grpcserver

import (
	"backend/backend/generated/auth"
	v1 "backend/handlers/v1"
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	resp, statusCode, err := v1.LoginHandlerLogic(req)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return resp, nil
}
