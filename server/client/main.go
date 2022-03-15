package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	"fmt"
	"log"

	"google.golang.org/grpc"
)
func main(){
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	c,e:=grpc.Dial("localhost:8081",grpc.WithInsecure())
	if e!=nil{
		log.Fatalf("cannot connect:%v",e)
	}
	tsClient:=trippb.NewTripServiceClient(c)
	r,err:=tsClient.GetTrip(context.Background(),&trippb.GetTripRequest{Id:"trip345"})
   if err!=nil{
	   log.Fatalf("cannot get trip:%v",err)
   }
   fmt.Println(r)
}