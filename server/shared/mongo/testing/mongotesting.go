package mongotesting

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)
var mongoURI string
const defaultmongoURI="mongodb://localhost:27017"
func RunwithMongo(m *testing.M)int{
	c,err:=client.NewEnvClient()
	if err!=nil{
		panic(err)
	}
	ctx:=context.Background()
	resp,err:=c.ContainerCreate(ctx,&container.Config{
		Image:"mongo:latest",
		ExposedPorts:nat.PortSet{
			"27017/tcp":{},
			},
		},&container.HostConfig{
			PortBindings:nat.PortMap{
				"27017/tcp":[]nat.PortBinding{
					{
						HostIP:"127.0.0.1",
						HostPort:"0",
					},

				},
			},
		},nil,nil,"")
	if err!=nil{
		panic(err)
	}
	defer func(){
		err:=c.ContainerRemove(ctx,resp.ID,types.ContainerRemoveOptions{
			Force:true,
		})
		if err!=nil{
			panic(err)
		}
	}()
	err=c.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{})
	if err!=nil{
		panic(err)
	}
	// time.Sleep(5*time.Second)
	inspRes,err:=c.ContainerInspect(ctx,resp.ID);
	if err!=nil{
		panic(err)
	}
	hostPort:=inspRes.NetworkSettings.Ports["27017/tcp"][0]
	mongoURI=fmt.Sprintf("mongodb://%s:%s",hostPort.HostIP,hostPort.HostPort)
	return m.Run()
}
func NewClient(c context.Context)(*mongo.Client,error){
	if mongoURI==""{
		return nil,fmt.Errorf("mongo not set RunwithMongo in TestMain")
	}
	return mongo.Connect(c,options.Client().ApplyURI(mongoURI))
}

func DefaultNewClient(c context.Context)(*mongo.Client,error){

	return mongo.Connect(c,options.Client().ApplyURI(defaultmongoURI))
}

func SetupIndexes(c context.Context,d *mongo.Database)error{
	_,err:=d.Collection("account").Indexes().CreateOne(c,mongo.IndexModel{
		Keys:bson.D{
			{Key:"open_id",Value:1},

		},
		Options: options.Index().SetUnique(true),
	})
	_,err=d.Collection("trip").Indexes().CreateOne(c,mongo.IndexModel{
		Keys:bson.D{
			{Key:"trip.accountid",Value:1},
			{Key:"trip.status",Value:1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.M{
				"trip.status":1,
			}),
	})
	return err
}
