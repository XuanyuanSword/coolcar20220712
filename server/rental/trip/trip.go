package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

type Service struct {
	//需要使用的参数
	Logger *zap.Logger
	Mongo *dao.Mongo
	ProfileManage ProfileManager
	CarManage CarManager
	POIManage POIManager
	DistanceCalc DistanceCalc
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


type DistanceCalc interface{
	DistanceKM(ctx context.Context,from  *rentalpb.Location,to *rentalpb.Location)(float64,error)
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	s.Logger.Info("CreateTrip",zap.String("code",req.CarId))
	//TODO: 验证驾驶者身份
	//TODO: 车辆开锁
	//TODO: 创建行程 开始计费
	aid,err:=auth.AccountID(c)
	s.Logger.Info("CreateTrip",zap.String("aid",string(aid)))
    if err!=nil{
    	return nil,err
	}
    if req.CarId==""||req.Start==nil{
       return nil,status.Error(codes.InvalidArgument,"")
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

	ls :=s.calcCurrentStatus(c,&rentalpb.LocationStatus{
		Location: req.Start,
		TimestampSec: nowFunc(),

	},req.Start)
	//poi,err:=s.POIManage.Resolve(c,req.Start)
	//
	//if err!=nil{
	//	s.Logger.Info("cannot resolve poi,",zap.Stringer("location",req.Start),zap.Error(err))
	//}
	//
	//ls:=&rentalpb.LocationStatus{
	//	Location:req.Start,
	//	PoiName: poi,
	//}

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

	s.Logger.Info("err",zap.Error(err))
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
	aid,err:=auth.AccountID(c)
	s.Logger.Info("GetTrip",zap.String("aid",string(aid)))
	if err!=nil{
		return nil,err
	}
	tr,err:=s.Mongo.GetTrip(c,id.TripID(req.Id),aid)
	if err!=nil{
		return nil,status.Error(codes.NotFound,"")
	}

	return tr.Trip,nil
}
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	aid,err:=auth.AccountID(c)
	s.Logger.Info("GetTrip",zap.String("aid",string(aid)))
	if err!=nil{
		return nil,err
	}
	trs,err:=s.Mongo.GetTrips(c,aid,req.Status)
	if err!=nil{
		s.Logger.Error("cannot get trips",zap.Error(err))
        return nil,status.Error(codes.Internal,"")
	}

    res :=&rentalpb.GetTripsResponse{

	}
	for _,tr:=range trs{
		res.Trips=append(res.Trips,&rentalpb.TripEntity{
			Id:tr.ID.Hex(),
			Trip:tr.Trip,
		})
	}
	return res,nil
}
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripReq) (*rentalpb.Trip, error) {
	// TDDO:为什么这里能够取到aid
	aid,err:=auth.AccountID(c)
	if  err!=nil{
		return nil,status.Error(codes.Unauthenticated,"")
	}
	tid:=id.TripID(req.Id)
	tr,err:=s.Mongo.GetTrip(c,tid,aid)
	if err!=nil{
		return nil,status.Error(codes.NotFound,"")
	}
	if tr.Trip.Current==nil{
		s.Logger.Error("trip without current set",zap.String("id",tid.String()))
	}
	cur:=tr.Trip.Current.Location

	if req.Current!=nil{
		cur=req.Current
	}
	tr.Trip.Current=s.calcCurrentStatus(c,tr.Trip.Current,cur)
	if req.EndTrip{
		tr.Trip.End=tr.Trip.Current
		tr.Trip.Status=rentalpb.TripStatus_FINISHED
	}
	s.Mongo.UpdateTrip(c,id.TripID(req.Id),aid,time.Now().Unix(),tr.Trip)
	return nil, status.Error(codes.Unimplemented, "")
}
const centsPerSec=0.7
var nowFunc=func() int64{
	return time.Now().Unix()
}
func (s *Service) calcCurrentStatus(c context.Context,last *rentalpb.LocationStatus,cur *rentalpb.Location)*rentalpb.LocationStatus{
	now:=nowFunc()
	//float64 转成浮点数
    elasedSec:=float64(now-last.TimestampSec)
    dist,err := s.DistanceCalc.DistanceKM(c,last.Location,cur)
    if err!=nil{
			s.Logger.Warn("cannot calculate distance",zap.Error(err))
	}
	poi,err:=s.POIManage.Resolve(c,cur)
	if err!=nil{
		s.Logger.Info("cannot resolve poi,",zap.Stringer("location",cur))
	}
	return &rentalpb.LocationStatus{
		Location: cur,
		//int32 转成整数
		FeeCent:last.FeeCent+int32(centsPerSec*elasedSec*2*rand.Float64()),
		KmDriven:last.KmDriven+dist,
		TimestampSec: now,
		PoiName:poi,
	}
}