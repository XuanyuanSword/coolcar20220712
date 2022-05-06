package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/server"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

//gateway
func main() {
	lg, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}
	//创造上下文
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	//grpc-gateway的一个请求多路复用器
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true, //枚举类型转换为整数
			OrigName:    true, //原始名称
		},
	))
	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "localhost:8082",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}
	for _, v := range serverConfig {
		err := v.registerFunc(
			c, mux, v.addr,
			[]grpc.DialOption{grpc.WithInsecure()},
		)
		if err != nil {
			lg.Sugar().Fatalf("cannot regisiter %s service:%v", v.name, err)
		}
	}
	addr := ":8080"
	lg.Sugar().Infof("start http server on %s", addr)
	lg.Sugar().Fatal(http.ListenAndServe(addr, mux))
}
