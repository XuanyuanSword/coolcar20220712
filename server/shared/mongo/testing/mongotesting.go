package mongotesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)
func RunwithMongo(m *testing.M,mongoURI *string)int{
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
	*mongoURI=fmt.Sprintf("mongodb://%s:%s",hostPort.HostIP,hostPort.HostPort)
	return m.Run()
}