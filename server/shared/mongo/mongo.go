package Mgo

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//IDFieldName
const (
	IDField            = "_id"
	UpdatedAtFieldName = "updatedat"
)

//IDField
type ObjId struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdatedAtField struct {
	UpdateAt int64 `bson:"updatedat"`
}

var NewObjID=primitive.NewObjectID

func NewObjIDwithValue(id fmt.Stringer){
	NewObjID=func() primitive.ObjectID{
		return objid.MustFromID(id)
	}
}
var UpdateAt=func() int64{
	return time.Now().UnixNano()
}

func Set(v interface{}) bson.M {
	return bson.M{"$set": v}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{"$setOnInsert": v}
}
