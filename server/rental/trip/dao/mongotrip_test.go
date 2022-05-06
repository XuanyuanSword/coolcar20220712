
package dao
import (
"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	"coolcar/shared/mongo/objid"

	mongotesting "coolcar/shared/mongo/testing"
"os"
"testing"

"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)
var mongoURI string
func TestCreateTrip(t *testing.T) {
	c:=context.Background()

	mc,err:=mongo.Connect(c,options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	m:=NewMongo(mc.Database("coolcar"))
	acct:=id.AccountIDs("account1")
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
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude: 35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err!=nil{
		t.Errorf("cannot create trip:%v",err)
	}
	t.Errorf("inserted row %s with updatedat %v",trip.ID,trip.UpdateAt)
	got,err:=m.GetTrip(c,objid.ToTripID(trip.ID),acct)
	if err!=nil{
		t.Errorf("cannot get trip：%v",err)
	}
	t.Errorf("got trip :%+v",got)

}
func TestMain(m *testing.M){
	os.Exit(mongotesting.RunwithMongo(m,&mongoURI))
}
