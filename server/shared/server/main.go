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
	nameField:=zap.String("name",c.Name)
	lis,err:=net.Listen("tcp",c.Addr)
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
		fmt.Println("in:",in,grpc.UnaryInterceptor(in));
		opts=append(opts,grpc.UnaryInterceptor(in))
	}
	s:=grpc.NewServer(opts...)
	c.RegisterFunc(s)
	c.Logger.Info("start grpc server",nameField,zap.String("addr",c.Addr))
	return s.Serve(lis)
}