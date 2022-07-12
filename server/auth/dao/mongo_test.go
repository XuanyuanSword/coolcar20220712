package dao

import (
	"context"
	"coolcar/shared/id"
	Mgo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var mongoURI string
func TestResolvrAccountId(t *testing.T) {
	c:=context.Background()

	mc,err:=mongo.Connect(c,options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err!=nil{
		t.Fatalf("connect mongodb failed,err:%v",err)
	}
	m:=NewMongo(mc.Database("coolcar"))
	_,err=m.col.InsertMany(c,[]interface{}{
		bson.M{
			Mgo.IDField:objid.MustFromID(id.AccountIDs("623197217465ba8f1d85bcb0")),
			openIDField:"openid1",
		},
		bson.M{
			Mgo.IDField:objid.MustFromID(id.AccountIDs("623197217465ba8f1d85bcb1")),
			openIDField:"openid2",
		},
	})
	if err!=nil{
		t.Fatalf("insert many failed,err:%v",err)
	}

	Mgo.NewObjIDwithValue(id.AccountIDs("623197217465ba8f1d85bcb2"))
	cases := []struct{
		name string
		openId string
		want id.AccountIDs
	}{
		{
			name:"存在用户",
			openId: "openid1",
			want: "623197217465ba8f1d85bcb0",
		},
		{
			name:"存在用户2",
			openId: "openid2",
			want: "623197217465ba8f1d85bcb1",
		},
		{
			name:"新用户",
			openId: "openid3",
			want: "623197217465ba8f1d85bcb2",
		},
	}
	for _,cc:=range cases{
		t.Run(cc.name,func(t *testing.T) {
			id,err:=m.ResolveAccountID(context.Background(),cc.openId)
			if err!=nil{
				t.Errorf("resolve account id failed %q,err:%v",cc.openId,err)
			}else{
				
				if id!=cc.want{
					t.Errorf("resolve account id failed,want:%s,got:%s",cc.want,id)
				}else{
					t.Logf("resolve account id success,id:%s",id)
				}
			}
		})
	}

}
func TestMain(m *testing.M){
	os.Exit(mongotesting.RunwithMongo(m))
}
//func mushObjId(hex string)primitive.ObjectID{
//	objID,_:=primitive.ObjectIDFromHex(hex)
//	return objID
//}