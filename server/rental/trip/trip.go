package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	//需要使用的参数
	Logger *zap.Logger
	Mongo *dao.Mongo
	ProfileManage ProfileManager
	CarManage CarManager
	POIManage POIManager
}
// 防腐层 （Anti Corruption Layer） 验证身份
type ProfileManager interface{
	Verify(context.Context,id.AccountIDs)(id.IdentityID,error)
}
//Resolve
type POIManager interface{
	Resolve(context.Context,*rentalpb.Location)(string,error)
}
//车辆管理
type CarManager interface{
	Verify(context.Context,id.CarID,*rentalpb.Location)(error)
	Unlock(context.Context,id.CarID)(error)
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	//TODO: 验证驾驶者身份
	//TODO: 车辆开锁
	//TODO: 创建行程 开始计费
	aid,err:=auth.AccountID(c)
    if err!=nil{
    	return nil,err
	}
	iID,err:=s.ProfileManage.Verify(c,aid)
	if err!=nil{
		return nil,status.Error(codes.FailedPrecondition,err.Error())
	}
	carID:=id.CarID(req.CarId)
	err=s.CarManage.Verify(c,carID,req.Start)
	if err!=nil{
		return nil,status.Error(codes.FailedPrecondition,err.Error())
	}
	poi,err:=s.POIManage.Resolve(c,req.Start)

	if err!=nil{
		s.Logger.Info("cannot resolve poi,",zap.Stringer("location",req.Start),zap.Error(err))
	}

	ls:=&rentalpb.LocationStatus{
		Location:req.Start,
		PoiName: poi,
	}

	tr,err:=s.Mongo.CreateTrip(c,&rentalpb.Trip{
		AccountId: aid.String(),
		CarId: carID.String(),
		IdentityId: iID.String(),
		Status: rentalpb.TripStatus_IN_PROGRESS,
		Start:ls,
		Current:ls,
	})
	if err!=nil{
		s.Logger.Warn("无法创建trip",zap.Error(err))
		return nil,status.Error(codes.AlreadyExists,"")
	}
	go func(){
		err=s.CarManage.Unlock(c,carID)
		if err!=nil{
			s.Logger.Error("无法开锁")
		}
	}()
	return &rentalpb.TripEntity{
		Id:tr.ID.Hex(),
		Trip:tr.Trip,
	},nil
}
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {

	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {

	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripReq) (*rentalpb.Trip, error) {
	// TDDO:为什么这里能够取到aid
	aid,err:=auth.AccountID(c)
	if  err!=nil{
		return nil,status.Error(codes.Unauthenticated,"")
	}
	tr,err:=s.Mongo.GetTrip(c,id.TripID(req.Id),aid)
	if req.Current!=nil{
		tr.Trip.Current=s.calcCurrentStatus(tr.Trip,req.Current)

	}
	if req.EndTrip{
		tr.Trip.End=tr.Trip.Current
		tr.Trip.Status=rentalpb.TripStatus_FINISHED
	}
	s.Mongo.UpdateTrip(c,id.TripID(req.Id),aid,time.Now().Unix(),tr.Trip)
	return nil, status.Error(codes.Unimplemented, "")
}
func (s *Service) calcCurrentStatus(trip *rentalpb.Trip,location *rentalpb.Location)*rentalpb.LocationStatus{
	return nil
}