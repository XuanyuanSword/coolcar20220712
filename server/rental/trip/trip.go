package auth

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Service struct{
	//需要使用的参数
	Logger *zap.Logger

}

func (s *Service) CreateTrip(c context.Context,req *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error){
	aid,err:=auth.AccountID(c)
	if err!=nil{
		s.Logger.Error("CreateTrip",zap.Error(err))
		fmt.Println(aid)
		return nil,err
	}
	return nil,status.Error(codes.Unimplemented,"")
}