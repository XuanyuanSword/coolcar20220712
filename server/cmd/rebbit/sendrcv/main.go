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

	q,err:=ch.QueueDeclare(
		"go_q1",
		true,
		false,
		false,
		false,
		nil,

		)
	if err!=nil{
		panic(err)
	}

	go consume("c1",conn,q.Name)
	go consume("c2",conn,q.Name)
	i:=0
	for{
		i++
		err:=ch.Publish(
			"",
			 q.Name,
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



func consume(name string,conn *amqp.Connection,q string){
     ch,err:=conn.Channel()
     if err!=nil{
     	panic(err)
	 }
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
