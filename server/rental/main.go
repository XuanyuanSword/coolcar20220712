package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	trip "coolcar/rental/trip"

	"coolcar/shared/server"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)
func main(){
	logger,err:=server.NewZapLogger()
	if err!=nil{
		log.Fatalf("cannot create logger :%v",err)
	}
	err=server.RunGRPCServer(&server.GRPCConfig{
		Name: "rental",
		Addr: ":8082",
		AuthPublicKeyFile: "shared/auth/public.key",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s,&trip.Service{
				Logger:logger,
			})
		},
	})
	// lis,err:=net.Listen("tcp",":8082")
	// if err!=nil{
	// 	logger.Fatal("cannot listen",zap.Error(err))
	// }
	// in,err:=auth.Interceptor("shared/auth/public.key");
	// if err!=nil{
	// 	logger.Fatal("cannot create interceptor",zap.Error(err))
	// }
	// fmt.Println("in:",in);
	// s:=grpc.NewServer(grpc.UnaryInterceptor(in))
	
	// rentalpb.RegisterTripServiceServer(s,&trip.Service{
	// 	Logger:logger,
	// })
	// err=s.Serve(lis)
	logger.Fatal("cannot server",zap.Error(err))
}

