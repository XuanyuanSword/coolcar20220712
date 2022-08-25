package main

import (
	"context"
	"coolcar/rental/ai"
	rentalpb "coolcar/rental/api/gen/v1"
	trip "coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/client/profile"
	"coolcar/rental/trip/dao"
	coolenvpb "coolcar/shared/coolenv"
	"coolcar/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}
	//62b9086d85a5f907a078047f
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err != nil {
		logger.Fatal("cannot connect mongo", zap.Error(err))
	}
	ac,err:=grpc.Dial("localhost:18001",grpc.WithInsecure())
	if err!=nil{
		logger.Fatal("cannot connect aisevice",zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "../shared/auth/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
				CarManage: &car.Manager{},
				ProfileManage: &profile.Manager{},
				POIManage: &poi.Manager{},
				DistanceCalc: &ai.Client{
					AIClient:coolenvpb.NewAIServiceClient(ac),
				},
				Mongo: dao.NewMongo(mongoClient.Database("coolcar")),
				})

		},
	}))

	logger.Fatal("cannot server", zap.Error(err))
}
