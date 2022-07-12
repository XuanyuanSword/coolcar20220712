package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T){
      //c:=context.Background()
	 c:=auth.ContextWithAccountId(context.Background(),"account1")
      mc,err:=mongotesting.NewClient(c)
      if err!=nil{
      	t.Fatalf("无法创建mongo 客户端:%v",err)
	  }
	  logger,err:=server.NewZapLogger()
	  if err!=nil{
	  	t.Fatalf("创建日志失败：%v",err)
	  }

      pm:=&profileManager{}
	  cm:=&carManager{}
	  s:=&Service{
	  	ProfileManage: pm,
	  	CarManage: cm,
	  	POIManage: &poi.Manager{},
	  	Mongo:dao.NewMongo(mc.Database("coolcar")),
	  	Logger: logger,
	  }
	  req:=&rentalpb.CreateTripRequest{
	  	CarId: "car1",
	  	Start:&rentalpb.Location{
	  		Latitude: 32.123,
	  		Longitude: 114.2525,
		},
	  }
	  pm.iID="identity1"
	  golden := `{"account_id":%q,"car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1605695246},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1605695246},"status":1,"identity_id":"identity1"}`
	  cases:=[]struct{
	  	name string
	  	tripID string
	  	profileErr error
	  	carVerifyErr error
	  	carUnlockErr error
	  	want string
	  	wantErr bool
	  }{
	  	{
	  		name:"normal_create",
	  		tripID:"6276399a91dcb29d940cfbb6",
	  		want:golden,//???

		},
		{
			  name:"profile_err",
			  tripID:"6276399a91dcb29d940cfbb7",
			  profileErr: fmt.Errorf("profile"),
			  wantErr:true,

		  },
		  {
	  		name:"car_verify_err",
	  		tripID: "6276399a91dcb29d940cfbb8",
	  		carVerifyErr: fmt.Errorf("verify"),
	  		wantErr: true,
		  },
		  {
			  name:"car_unlock_err",
			  tripID: "6276399a91dcb29d940cfbb9",
			  carVerifyErr: fmt.Errorf("unlock"),
			  want: golden,
		  },
	  }
	  for _,cc:=range cases{
	  	t.Run(cc.name,func(t *testing.T){
			Mgo.NewObjIDwithValue(id.TripID(cc.tripID))
			pm.err=cc.profileErr
			cm.unlockErr=cc.carUnlockErr
			cm.verifyErr=cc.carVerifyErr

			res,err:=s.CreateTrip(c,req)
			if cc.wantErr{
				if err==nil{
					t.Errorf("want error;got none")
				}else{
					return
				}
			}
			if err!=nil{
				t.Errorf("无效 创建trip %v",err)
				return
			}

			if res.Id!=cc.tripID{
				t.Errorf("无效id； want %q, got %q",cc.tripID,res.Id)
			}

			b,err:=json.Marshal(res.Trip)
			if err!=nil{
				t.Errorf("无法 格式化 响应:%v",err)
			}

			got:=string(b)
			if cc.want !=got{
				t.Errorf("无效响应:want %q,got %q",cc.want,got)
			}


		})
	  }
}

type profileManager struct{
     iID id.IdentityID
     err error

}

func (p *profileManager) Verify(context.Context,id.AccountIDs)(id.IdentityID,error){
	return p.iID,p.err
}

type carManager struct{
	verifyErr error
	unlockErr error
}

func (c *carManager) Verify(context.Context,id.CarID,*rentalpb.Location)(error){
    return c.verifyErr
}

func (c *carManager) Unlock(context.Context,id.CarID)(error){
	return c.unlockErr
}


func TestMain(m *testing.M){
	os.Exit(mongotesting.RunwithMongo(m))
}