package dao

import (
	"context"
	Rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//const openIDField = "openid"
const (
	tripField = "trip"
	accountIDField=tripField+".accountid"

)
type Mongo struct {
	col      *mongo.Collection
	newObjID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("trip"),
		newObjID: primitive.NewObjectID,
	}
}

type TripRecord struct {
	Mgo.ObjId `bson:"inline"`
	Mgo.UpdatedAtField `bson:"inline"`
	Trip     *Rentalpb.Trip `bson:"trip"`
}
// TODO: 同一个account最多只能又一个进行中的Trip
// TODO: 强类型化tripID
// TODO: 表格驱动测试
func (m *Mongo) CreateTrip(c context.Context, trip *Rentalpb.Trip) (*TripRecord, error) {
	// var t TripRecord
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = m.newObjID()
	r.UpdateAt = time.Now().UnixNano()
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
