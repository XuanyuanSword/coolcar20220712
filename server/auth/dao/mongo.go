package dao

import (
	"context"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const openIDField="openid";

type Mongo struct {
	col *mongo.Collection

}
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:db.Collection("account"),

	}
}

func (m *Mongo) ResolveAccountID(c context.Context,openID string) (id.AccountIDs, error) {
	// m.col.InsertOne(c,bson.M{
	// 	mgo.IDField:m.newObjID(),
	// 	openIDField:openID,
	// })
	insertID:=Mgo.NewObjID()
	res := m.col.FindOneAndUpdate(c,bson.M{
		openIDField:openID,
	},
	Mgo.SetOnInsert(bson.M{
			Mgo.IDField:insertID,
			openIDField:openID,
		}), 
	options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After))
	if err:=res.Err();err!=nil{
		return "",fmt.Errorf("resolve account id failed,err:%v",err)
	}
	var row  Mgo.ObjId
	err:=res.Decode(&row)
	if err!=nil{
		return "",fmt.Errorf("resolve account id failed,err:%v",err)
	}
	return objid.ToAccountID(row.ID),nil
}