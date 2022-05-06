package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"coolcar/secret"
	"coolcar/shared/server"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewDevelopment()
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err != nil {
		logger.Fatal("cannot connect mongo", zap.Error(err))
	}
	tk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret.PrivateKey))
	if err != nil {
		//测试会中止
		logger.Fatal("Failed to parse RSA private key", zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name: "auth",
		Addr: ":8081",
		// AuthPublicKeyFile: "shared/auth/public.key",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID:     "wx3e974e6b62b3b907",
					AppSecret: secret.AppSecret,
				},
				Logger:         logger,
				TokenGenerator: token.NewJWTToken("coolcar/auth", tk),
				TokenExpire:    2 * time.Hour,
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
			})
		},
	}))
}
