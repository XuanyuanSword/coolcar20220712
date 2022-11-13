package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)
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
