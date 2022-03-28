package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Service struct{
	//需要使用的参数
	Logger *zap.Logger
	OpenIDResolver OpenIDResolverType
	TokenGenerator Tokengen
	Mongo *dao.Mongo
	TokenExpire time.Duration
}
type OpenIDResolverType interface {
	//参数中 需要定义的方法 wechat 
	Resolve(code string) (string,error)
}
type Tokengen interface{
	 GenerateToken(accountID string,expire time.Duration) (string,error)
}
func (s *Service) Login(c context.Context,req *authpb.LoginRequest) (*authpb.LoginResponse, error){
	s.Logger.Info("Login",zap.String("code",req.Code))
	openid,err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil {
		s.Logger.Error("Login",zap.Error(err))
		return nil,status.Errorf(codes.Unavailable,"不能返回openid%v",err)
	}
	accID,err:=s.Mongo.ResolveAccountID(c,openid)

	if err != nil {
		s.Logger.Error("Login",zap.Error(err))
		return nil,status.Errorf(codes.Unavailable,"不能返回accid%v",err)
	}
	tkn,err:=s.TokenGenerator.GenerateToken(accID,s.TokenExpire)
	if err != nil {
		s.Logger.Error("Login",zap.Error(err))
		return nil,status.Errorf(codes.Unavailable,"不能生成token%v",err)
	}
	return &authpb.LoginResponse{
		AccessToken: "token accID"+tkn,
		ExpiresIn:  int32(s.TokenExpire.Seconds()),
	},nil
}