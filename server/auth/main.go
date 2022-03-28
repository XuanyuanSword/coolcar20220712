package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"log"
	"net"
	"time"

	"coolcar/secret"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)
func main(){
	logger,err:=zap.NewDevelopment()
	if err!=nil{
		log.Fatalf("cannot create logger :%v",err)
	}
	lis,err:=net.Listen("tcp",":8081")
	if err!=nil{
		logger.Fatal("cannot listen",zap.Error(err))
	}
	c:=context.Background()
	mongoClient,err:=mongo.Connect(c,options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err!=nil{
		logger.Fatal("cannot connect mongo",zap.Error(err))
	}
	s:=grpc.NewServer()
	tk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret.PrivateKey))
	if err != nil {
		//测试会中止
		logger.Fatal("Failed to parse RSA private key",zap.Error(err))
	}
	authpb.RegisterAuthServiceServer(s,&auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID: "wx3e974e6b62b3b907",
			AppSecret: secret.AppSecret,

		},
		Logger:logger,
		TokenGenerator:token.NewJWTToken("coolcar/auth",tk),
		TokenExpire: 2*time.Hour,
		Mongo:dao.NewMongo(mongoClient.Database("coolcar")),
	})

   err=s.Serve(lis)
   logger.Fatal("cannot server",zap.Error(err))
}

func newZapLogger()(*zap.Logger,error){
	logger:=zap.NewDevelopmentConfig()
	logger.EncoderConfig.TimeKey=""
	return logger.Build()
}