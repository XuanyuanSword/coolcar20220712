package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func main(){
	c:=context.Background()
	mc,err:=mongo.Connect(c,
	options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err!=nil{
		panic(err)
	}
	col:=mc.Database("coolcar").Collection("account")
	res,err:=col.InsertMany(c,[]interface{}{
		bson.M{"name":"张三","age":18},
		bson.M{"name":"李四","age":19},
	})
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%+v",res)
	findRow(c,col)
}
func findRow(c context.Context,col *mongo.Collection) {
	res,err :=col.Find(c,bson.M{})
	if err!=nil{
		panic(err)
	}
	for res.Next(c){
		var row struct {
			ID primitive.ObjectID `bson:"_id"`
			Name string `bson:"name"`
		}
		if err=res.Decode(&row);err!=nil{
			panic(err)
		}
		fmt.Printf("%+v\n",row)
	}
}
func insertRow(c context.Context,col *mongo.Collection) {
	res,err:=col.InsertMany(c,[]interface{}{
		bson.M{"name":"张三","age":18},
		bson.M{"name":"李四","age":19},
	})
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%+v",res)
}