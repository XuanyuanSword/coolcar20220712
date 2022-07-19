package dao

import (
	"context"
	Rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//const openIDField = "openid"
const (
	tripField = "trip"
	accountIDField=tripField+".accountid"
	statusField    = tripField + ".status"
)
type Mongo struct {
	col      *mongo.Collection

}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("trip"),

	}
}

type TripRecord struct {
	Mgo.ObjId `bson:"inline"`
	Mgo.UpdatedAtField `bson:"inline"`
	Trip     *Rentalpb.Trip `bson:"trip"`
}
// TODO: 同一个account最多只能又一个进行中的Trip    建monggo 索引 =》 shared/mongo/set.js
// TODO: 强类型化tripID => shared/id/id.go,shared/mongo/objid/objid.go
// TODO: 表格驱动测试 =>
func (m *Mongo) CreateTrip(c context.Context, trip *Rentalpb.Trip) (*TripRecord, error) {
	// var t TripRecord
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = Mgo.NewObjID()

	r.UpdateAt = Mgo.UpdateAt()
	_,err:=m.col.InsertOne(c,r)
	if err!=nil{
		return nil,err
	}
	return r,nil

}
// getTrip
func (m *Mongo)GetTrip(c context.Context,id  id.TripID,accountID id.AccountIDs)(*TripRecord,error){
	//string 转 ObjectID
	objID,err:=objid.FromID(id)

	if err!=nil{
		return nil,fmt.Errorf("invaild id: %v",err)
	}
	res:=m.col.FindOne(c,bson.M{
		Mgo.IDField: objID,
		accountIDField:accountID,
	})
	if err:=res.Err();err!=nil{
		return nil,err
	}

	var tr TripRecord
	err=res.Decode(&tr)
	if err!=nil{
		return nil,fmt.Errorf("cannot decodeL%v",err)
	}
	return &tr,nil
}

func (m *Mongo)GetTrips(c context.Context,accountID id.AccountIDs,status Rentalpb.TripStatus)([]*TripRecord,error){
	filter:=bson.M{
		accountIDField: accountID.String(),
	}
	if status!=Rentalpb.TripStatus_TS_NOT_SPECIFIED{
		filter[statusField]=status;
	}
	res,err:=m.col.Find(c,filter)
	if err!=nil{
		return nil,err
	}
	var trips []*TripRecord
	for res.Next(c){
		var trip TripRecord
		err:=res.Decode(&trip)
		if err!=nil{
			return nil,err
		}
		trips=append(trips,&trip)
	}
	return trips,nil
}

func (m *Mongo)UpdateTrip(c context.Context,tid id.TripID,aid id.AccountIDs,updateAt int64,trip *Rentalpb.Trip)error{
	objID,err:=objid.FromID(tid)
	if err!=nil{
		return fmt.Errorf("无效id %v",err)
	}
	newUpdatedAt:=Mgo.UpdateAt()
	res,err:=m.col.UpdateOne(c,bson.M{
		Mgo.IDField: objID,
		accountIDField: aid.String(),
		Mgo.UpdatedAtFieldName: updateAt,

	},Mgo.Set(bson.M{
		tripField: trip,
		Mgo.UpdatedAtFieldName: newUpdatedAt,
	}))
	if err!=nil{
		return err
	}
	if res.MatchedCount==0{
		return mongo.ErrNoDocuments
	}
	return err
}