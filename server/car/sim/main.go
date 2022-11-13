package sim

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"go.uber.org/zap"
)
 //收消息
type Subscriber interface {
	Subscribe(ctx context.Context)(chan *carpb.CarEntity,func(),error)
}
type Controller struct{
	CarService carpb.CarServiceClient
	Logger *zap.Logger
	Subscriber Subscriber
}
//type Subscribers struct{
//
//}
//func (Subscribers) Subscriber(ctx context.Context)(chan *carpb.CarEntity,error) {
//	ch := make(chan *carpb.CarEntity)
//	return ch,nil
//}

func (c *Controller) RunSimulations(ctx context.Context) {
	res, err := c.CarService.GetCars(ctx, &carpb.GetCarsRequest{})
	if err != nil {
		c.Logger.Error("cannot get cars ", zap.Error(err))
		return
	}
	msgCh,cleanUp,err:=c.Subscriber.Subscribe(ctx)
	defer  cleanUp()
	if err!=nil{
		c.Logger.Error("cannot subscribe",zap.Error(err))
		return
	}
    // 键类型 string 值类型 chan *carpb.Car
	carChans := make(map[string]chan *carpb.Car)
	//range遍历数组
	for _, car := range res.Cars {
		ch := make(chan *carpb.Car)
		carChans[car.Id] = ch
		go c.SimulateCar(context.Background(),car,ch)
	}
	//range遍历channel
	for carUpdate:=range msgCh{
		ch:=carChans[carUpdate.Id]
			if ch!=nil{
				ch<-carUpdate.Car
			}
	}
}
func (c *Controller) SimulateCar(ctx context.Context,initial *carpb.CarEntity,ch chan *carpb.Car){
	carID:=initial.Id
	for update:=range ch{
		if update.Status==carpb.CarStatus_UNLOCKING{
			_,err:=c.CarService.UpdateCar(ctx,&carpb.UpdateCarRequest{
				Id:carID,
				Status: carpb.CarStatus_UNLOCKED,
			})
			if err!=nil{
				c.Logger.Error("cannot unlock",zap.Error(err))
			}
		} else if update.Status==carpb.CarStatus_LOCKING{
			_,err:=c.CarService.UpdateCar(ctx,&carpb.UpdateCarRequest{
				Id:carID,
				Status: carpb.CarStatus_LOCKED,
			})
			if err!=nil{
				c.Logger.Error("cannot lock",zap.Error(err))
			}
		}
	}

}