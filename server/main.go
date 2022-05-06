package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	// "google.golang.org/protobuf/proto"
)
func main(){
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	// startGRPCGateway()
	go startGRPCGateway()

	lis,err:=net.Listen("tcp",":8081")
	if err!=nil{
		log.Fatalf("failed to listen:%v",err)

		
	}
	s:=grpc.NewServer()
	trippb.RegisterTripServiceServer(s,&trip.Service{})
	log.Fatal(s.Serve(lis))
	// //
	// var a int
	// fmt.Println(a)
	// trip:=trippb.Trip{
	// 	Start:"abc",
	// 	End:"def",
	// 	DurationSec:3600,
	// 	FeeCent:10000,
	// }
	// fmt.Println("Hello World",&trip)
    // b,err:=proto.Marshal(&trip)
	// if err!=nil{
	// 	panic(err)
	// }
	// fmt.Printf("%X\n",b)
	// var trip2 trippb.Trip
	// err=proto.Unmarshal(b,&trip2)
	// if err!=nil{
	// 	panic(err)
	// }
	// fmt.Println(&trip2)
	// b,err=json.Marshal(&trip)
	// if err!=nil{
	// 	panic(err)
	// }
	// fmt.Printf("%s\n",b)
}

func startGRPCGateway(){
	c:=context.Background()
	c,cancel:=context.WithCancel(c)
	defer cancel()
	mux:=runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,&runtime.JSONPb{
			EnumsAsInts:true,//枚举类型转换为整数
			OrigName: true,//原始名称
		},
	))
	err:=trippb.RegisterTripServiceHandlerFromEndpoint(c,
		//mux:multiplexer
		mux,
         ":8081",
		 []grpc.DialOption{grpc.WithInsecure()},
		)
	if err!=nil{
		log.Fatalf("failed to register gateway:%v",err)
	}
    err= http.ListenAndServe(":8080",mux)
	if err!=nil{
		log.Fatalf("failed to start gateway:%v",err)
	}
}
