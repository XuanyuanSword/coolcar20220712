package auth

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	//需要使用的参数
	Logger *zap.Logger
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	// aid,err:=auth.AccountID(c)
	// if err!=nil{
	// 	s.Logger.Error("CreateTrip",zap.Error(err))
	// 	fmt.Println("123",aid)
	// 	return nil,err
	// }
	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {

	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {

	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripReq) (*rentalpb.Trip, error) {

	return nil, status.Error(codes.Unimplemented, "")
}
