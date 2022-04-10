package server

import (
	"coolcar/shared/auth"
	"fmt"

	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)
type GRPCConfig struct{
	Name string
	Addr string
	AuthPublicKeyFile string
	Logger *zap.Logger
	RegisterFunc func(*grpc.Server)

}
func RunGRPCServer(c *GRPCConfig) error{
	// logger,err:=zap.NewDevelopment()
	// if err!=nil{
	// 	log.Fatalf("cannot create logger :%v",err)
	// }
	nameField:=zap.String("name",c.Name)
	lis,err:=net.Listen("tcp",":8082")
	if err!=nil{
		c.Logger.Fatal("cannot listen",nameField,zap.Error(err))
	}
	//"shared/auth/public.key"
	var opts []grpc.ServerOption
	if(c.AuthPublicKeyFile!=""){

		in,err:=auth.Interceptor(c.AuthPublicKeyFile);
		if err!=nil{
			c.Logger.Fatal("cannot create interceptor",zap.Error(err))
		}
		fmt.Println("in:",in);
		opts=append(opts,grpc.UnaryInterceptor(in))
	}
	s:=grpc.NewServer(opts...)
	c.RegisterFunc(s)
	return s.Serve(lis)
	// c.Re
	// rentalpb.RegisterTripServiceServer(s,&trip.Service{
	// 	Logger:c.Logger,
	// })
	// err=s.Serve(lis)
	// logger.Fatal("cannot server",zap.Error(err))
}