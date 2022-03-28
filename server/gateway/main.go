package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

//gateway
func main(){
	//创造上下文
	c:=context.Background()
	c,cancel:=context.WithCancel(c)
	defer cancel()
	//grpc-gateway的一个请求多路复用器
	mux:=runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,&runtime.JSONPb{
			EnumsAsInts:true,//枚举类型转换为整数
			OrigName: true,//原始名称
		},
	))
	err:=authpb.RegisterAuthServiceHandlerFromEndpoint(
		c,mux,"localhost:8081",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err!=nil{
		log.Fatalf("cannot regisiter auth service:%v",err)
	}
	log.Fatal(http.ListenAndServe(":8080",mux))
}