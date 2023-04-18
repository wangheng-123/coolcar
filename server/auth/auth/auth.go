package auth

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authpb "server/auth/api/gen/v1"
	"server/auth/dao"
)

// Service implements auth Service
type Service struct {
	OpenIDResolver OpenIDResolver
	Logger         *zap.Logger
	Mongo          *dao.Mongo
}

type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// Login 用户登录服务
//Login logs a user in
func (s *Service) Login(c context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", request.Code))
	openID, err := s.OpenIDResolver.Resolve(request.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid；%v", err)
	}

	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: "token for account id：" + accountID,
		ExpiresIn:   7200,
	}, nil
}
