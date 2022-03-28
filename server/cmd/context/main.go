package main

import (
	"context"
	"fmt"
	"time"
)

func main(){
	c:=context.WithValue(context.Background(),"key","value")
	c,cancel:=context.WithTimeout(c,5*time.Second)
	
	defer cancel()
	mainTask(c)
	time.Sleep(time.Minute )
}

func mainTask(c context.Context){
	// fmt.Printf("%q\n",c.Value("key"))
	//退出函数时不执行 cancel 方法 c1可以执行完
	go func(){
		//子任务  context.Background() 新任务
		c1,cancel:=context.WithTimeout(c,10*time.Second)
		go smallTask(c1,"World1",9*time.Second)
		defer cancel()
	}()
	// smallTask(c,"World1",4*time.Second)
	smallTask(c,"World2",2*time.Second)
}

func smallTask(c context.Context, name string,d time.Duration){

	
	fmt.Println("start", name,c.Value("key"))
	select{
		case <-time.After(d):
			fmt.Println("done", name)
		case <-c.Done():
			fmt.Printf("%s is canceled.\n", name)
	}
	<-c.Done()
}
