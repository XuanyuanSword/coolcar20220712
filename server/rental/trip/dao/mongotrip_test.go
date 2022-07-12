
package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	"testing"
)
var mongoURI string

func TestCreateTrip(t *testing.T){
	c:=context.Background()

	mc,err:=mongotesting.NewClient(c)
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	db:=mc.Database("coolcar")
	err=mongotesting.SetupIndexes(c,db)
	if err!=nil{
		t.Fatalf("cannot setup indexes:%v",err)
	}
	m:=NewMongo(db)
	cases:=[]struct{
		name string
		tripID string
		accountID string
		tripStatus rentalpb.TripStatus
		wantErr bool
	}{
		{
			name: "完成",
			tripID: "6276399a91dcb29d940cfbb6",
			accountID: "account111",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name: "其他完成",
			tripID: "6276399a91dcb29d940cfbb7",
			accountID: "account111",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name: "in_process",
			tripID: "6276399a91dcb29d940cfbb8",
			accountID: "account111",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name: "其他in_process",
			tripID: "6276399a91dcb29d940cfbb9",
			accountID: "account111",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:true,
		},
		{
			name: "in_process不同accountID",
			tripID: "6276399a91dcb29d940cfb10",
			accountID: "account112",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,

		},
	}
	for _,cc:=range cases{
		Mgo.NewObjID=func() primitive.ObjectID{
			return objid.MustFromID(id.TripID(cc.tripID))
		}
		tr,err:=m.CreateTrip(c,&rentalpb.Trip{
			AccountId: cc.accountID,
			Status: cc.tripStatus,
		})
		if cc.wantErr{
			//没有错
			if err==nil{
				t.Errorf("===%s: error expected; got none",cc.name)

			}
			//直接走下个循环
			continue
		}
		//有错
		if err!=nil{
			t.Errorf("%s error createing trip:%v",cc.name,err)
			continue
		}
		if tr.ID.Hex()!=cc.tripID{
			t.Errorf("incorrect trip id;want;%q;got:%q",
				cc.tripID,tr.ID.Hex())
		}


	}
}

func TestGetTrip(t *testing.T) {
	c:=context.Background()

	mc,err:=mongotesting.NewClient(c)
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	m:=NewMongo(mc.Database("coolcar"))
	acct:=id.AccountIDs("account488")
	//这里的newObjID 不太一样
	Mgo.NewObjID=primitive.NewObjectID
	trip,err:=m.CreateTrip(c,&rentalpb.Trip{
		AccountId: acct.String(),
		CarId: "car1",
		Start:&rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude: 30,
				Longitude: 120,
			},
		},
		End:&rentalpb.LocationStatus{
			PoiName: "endpoint",
			FeeCent: 10000,
			KmDriven: 350,
			Location: &rentalpb.Location{
				Latitude: 35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err!=nil{
		t.Fatalf("cannot create trip:%v",err)
	}
	//t.Errorf("inserted row %s with updatedat %v",trip.ID,trip.UpdateAt)
	got,err:=m.GetTrip(c,objid.ToTripID(trip.ID),acct)

	if err!=nil{
		t.Errorf("cannot get trip：%v",err)
	}
	if diff:=cmp.Diff(trip,got,protocmp.Transform());diff!=""{
		t.Errorf("result differs; -want +got:%s",diff)
	}

}
func TestGetTrips(t *testing.T){
	rows:=[]struct{
		id string
		accountID string
		status rentalpb.TripStatus
	}{
		{

			id: "5276399a91dcb29d940cfbb6",
			accountID: "account_id_for_get_trips",
			status: rentalpb.TripStatus_FINISHED,
		},
		{
			id: "5276399a91dcb29d940cfbb7",
			accountID: "account_id_for_get_trips",
			status: rentalpb.TripStatus_FINISHED,
		},
		{
			id: "5276399a91dcb29d940cfbb8",
			accountID: "account_id_for_get_trips",
			status: rentalpb.TripStatus_FINISHED,
		},
		{
			id: "5276399a91dcb29d940cfbb9",
			accountID: "account_id_for_get_trips",
			status: rentalpb.TripStatus_IN_PROGRESS,

		},
		{
			id: "5276399a91dcb29d940cfb10",
			accountID: "account_id_for_get_trips_other",
			status: rentalpb.TripStatus_IN_PROGRESS,

		},

	}
	c:=context.Background()

	mc,err:=mongotesting.NewClient(c)
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	m:=NewMongo(mc.Database("coolcar"))
	for _,cc :=range rows{
		Mgo.NewObjIDwithValue(id.TripID(cc.id))
		_,err:=m.CreateTrip(c,&rentalpb.Trip{
			AccountId: cc.accountID,
			Status: cc.status,
		})
		if err !=nil{
			t.Fatalf("无法创建行程:%v",err)
		}

	}
	cases:=[]struct{
		name string
		accountID string
		status rentalpb.TripStatus
		wantCount int
		wantOnlyID string
	}{
		{
			name:"get_all",
			accountID: "account_id_for_get_trips",
			status:rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},{
			name:"get_in_progress",
			accountID: "account_id_for_get_trips",
			status: rentalpb.TripStatus_IN_PROGRESS,
			wantCount: 1,
			wantOnlyID: "5276399a91dcb29d940cfbb9",
		},

	}
	for _,cc:=range cases{
		t.Run(cc.name, func(t *testing.T) {
			res,err:=m.GetTrips(context.Background(),id.AccountIDs(cc.accountID),cc.status)
			if err!=nil{
				t.Errorf("获取行程错误%v",err)
			}
			if cc.wantCount!=len(res){
				t.Errorf("获取行程数不对，want %d ，got %d",cc.wantCount,len(res))
			}
			if cc.wantOnlyID!=""&&len(res)>0{
				if cc.wantOnlyID!=res[0].ID.Hex(){
					t.Errorf("行程的id不同 want %q,got %q",cc.wantOnlyID,res[0].ID.Hex())
				}
			}

		})
	}

}
func TestUpdateTrip(t *testing.T){
	c:=context.Background()

	mc,err:=mongotesting.NewClient(c)
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	m:=NewMongo(mc.Database("coolcar"))
	var now int64=10000
	tid:=id.TripID("6276399a91dcb29d940cfbb6")
	aid:=id.AccountIDs("update_id")
	Mgo.NewObjIDwithValue(tid)
	Mgo.UpdateAt= func() int64 {
		return now
	}
	tr,err:=m.CreateTrip(c,&rentalpb.Trip{
		AccountId: aid.String(),
		Status:rentalpb.TripStatus_IN_PROGRESS,
		Start:&rentalpb.LocationStatus{
			PoiName: "start_poi",
		},

	})
	if err!=nil{
		t.Fatalf("无法创建行程%v",err)
	}
	if tr.UpdateAt!=10000{
		t.Fatalf("want %d,got %d",now,tr.UpdateAt)
	}
	update:=&rentalpb.Trip{
		AccountId: aid.String(),
		Status:rentalpb.TripStatus_IN_PROGRESS,
		Start:&rentalpb.LocationStatus{
			PoiName: "start_poi_update",
		},
	}
	cases:= []struct{
		 name string
		 now int64
		 withUpdateAt int64
		 wantErr bool
	}{
		{
			name:"normal update",
			now:20000,//  这一次要更新的值
			withUpdateAt: 10000, //  上一次的值

		},
		//没有更新成功
		{
			name:"normal update error",
			now:30000,
			withUpdateAt: 10000,//shang
			wantErr: true,
		},
		//基于第一次更新成功
		{
			name:"update_with_refetch",
			now:40000,
			withUpdateAt: 20000,

		},
	}
	for _,cc:=range cases{
		now=cc.now
		err:=m.UpdateTrip(c,tid,aid,cc.withUpdateAt,update)
		if cc.wantErr{
			if err==nil{
				t.Errorf("%s :want error;got none",cc.name)
			}else{
				continue
			}
		}else{
			if err!=nil{
				t.Errorf("无法更新，%v",cc.name)
			}else{
				continue
			}
		}
		updateTrip,err:=m.GetTrip(c,tid,aid)
		if err!=nil{
			t.Errorf("更新之后无法获取行程%v",err)
		}
		if cc.now!=updateTrip.UpdateAt{
			t.Errorf("%s, 更新值不一致 want %d ，got %d",cc.name,cc.now,updateTrip.UpdateAt)
		}
	}

}
func TestMain(m *testing.M){
	os.Exit(mongotesting.RunwithMongo(m))
}
