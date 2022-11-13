package main
import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"encoding/json"
	"fmt"
"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)
type Publisher struct {
	ch *amqp.Channel
	exchange string
}
// 发消息
func NewPublisher(conn *amqp.Connection,exchange string)(*Publisher,error ) {

	//虚拟的链接 amqp的channel
	ch, err := conn.Channel()
	if err!=nil{
		return nil,fmt.Errorf("cannot allocate （分配） channel:%v",err)
	}

	err=declareExchange(ch,exchange)
	if err!=nil{
		return nil,fmt.Errorf("cannot declare（声明） exchange:%v",err)
	}

	return &Publisher{ch,exchange}, err
}
func (p *Publisher)Publisher(c context.Context,car *carpb.CarEntity)error{
	b,err:=json.Marshal(car)
	if err!=nil{
		return fmt.Errorf("cannot marshal:%v",err)
	}
	return p.ch.Publish(
		p.exchange,
		//"",
		//q.Name,
		"",
		false,
		false,
		amqp.Publishing{
			Body: b,
		},
	)

}


const exchangeName="go_ex"
func main(){
	conn,err:=amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		panic(err)
	}
	//虚拟的链接 amqp的channel
	ch,err:=conn.Channel()
	if err!=nil{
		panic(err)
	}
	err=ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)


	if err!=nil{
		panic(err)
	}
	go subscribe(conn,exchangeName)
	go subscribe(conn,exchangeName)
	i:=0
	for{
		i++
		err:=ch.Publish(
			exchangeName,
			//"",
			//q.Name,
			"",
			false,
			false,
			amqp.Publishing{
				Body: []byte(fmt.Sprintf("message %d",i)),
			},
		)
		if err!=nil{
			fmt.Println(err.Error())
		}
		time.Sleep(200*time.Millisecond)
	}


}

func subscribe(conn *amqp.Connection,ex string){
	ch,err:=conn.Channel()
	if err!=nil{
		panic(err)
	}
	defer ch.Close()
	q,err:=ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,

	)
	if err!=nil{
		panic(err)
	}
	defer ch.QueueDelete(
		q.Name,
		false,
		false,
		false,
	)
	ch.QueueBind(
		q.Name,
		"",
		ex,
		false,
		nil,
	)
	consume("c",ch,q.Name)
}

func consume(name string,ch *amqp.Channel,q string){

	msgs,err:=ch.Consume(
		q,
		name,
		true,
		false,
		false,
		false,
		nil,
	)

	if err!=nil{
		panic(err)
	}
	for msg:=range msgs{
		fmt.Printf("%s %s\n",name,msg.Body)
	}

}
type Subscriber struct{
	conn *amqp.Connection
	exchange string
	logger *zap.Logger
}
func (s *Subscriber)SubscribeRaw(ctx context.Context)(<-chan amqp.Delivery,func(),error){
	ch,err:=s.conn.Channel()
	if err!=nil{
		return nil,func(){},fmt.Errorf("cannot allocate channel:%v",err)
	}
	closeCh:= func() {
		err:=ch.Close()
		if err!=nil{
			// 这个日志 字符串不需要拼接参数
			s.logger.Error("cannot close channel",zap.Error(err))
		}

	}
	q,err:=ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,

	)
	if err!=nil{
		return nil,closeCh,fmt.Errorf("cannot declare queue: %v",err)
	}
	closeUp:= func() {
		_,err:=ch.QueueDelete(
			q.Name,
			false,
			false,
			false,
		)
		if err!=nil{
			s.logger.Error("cannot  declare queue: ",zap.String("name",q.Name))
		}
		closeCh()
	}

	//defer ch.QueueDelete(
	//	q.Name,
	//	false,
	//	false,
	//	false,
	//)
	ch.QueueBind(
		q.Name,
		"",
		s.exchange,
		false,
		nil,
	)
	if err!=nil{
		return nil,closeUp,fmt.Errorf("cannot bind:%v",err)
	}
	msgs,err:=ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err!=nil{
		return nil,closeUp,fmt.Errorf("")
	}

	return msgs,closeUp,nil
}
func (s *Subscriber) Subscribe(ctx context.Context)(chan *carpb.CarEntity,func(),error){
	 msgCh,cleanUp,err:=s.SubscribeRaw(ctx)
		if err!=nil{
			return nil,cleanUp,err
		}
	carCh:=make(chan *carpb.CarEntity)
	go func(){
		for msg:=range msgCh{
			var car carpb.CarEntity
			err:=json.Unmarshal(msg.Body,&car)
			if err!=nil{
				s.logger.Error("无法转译",zap.Error(err))
			}
			carCh <-&car

		}
		close(carCh)
	}()

	return carCh,cleanUp,nil
}
func NewSubscriber(conn *amqp.Connection,exchange string,logger *zap.Logger)(*Subscriber,error ){
	//虚拟的链接 amqp的channel
	ch, err := conn.Channel()
	if err!=nil{
		return nil,fmt.Errorf("cannot allocate （分配） channel:%v",err)
	}
	defer  ch.Close()
	err=declareExchange(ch,exchange)
	if err!=nil{
		return nil,fmt.Errorf("cannot declare（声明） exchange:%v",err)
	}

	return &Subscriber{conn,exchange,logger}, err
}
func declareExchange(ch *amqp.Channel,exchange string)error{
	return ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}