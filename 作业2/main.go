package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var done = make(chan error, 2)
	//外部利用stop对goroutine进行控制
	var stop = make(chan struct{})

	//程序结束了，得向done channel中发哦是那个
	go func() {
		done <- StartHTTPServer(stop)
	}()
	go func() {
		done <- HandleSingal(stop)
	}()

	//有两个goroutine等待stop，所以得关闭这个stop channel
	var stopped bool
	//接受两个服务的err
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			//如果done中有数据了，代表某个服务退出了，就要提醒另一个服务退出
			//stop <- struct{}{} 		//panic: send on closed channel
			fmt.Printf("err : %s \n", err.Error())
		}
		if !stopped { //防止关闭两次
			stopped = true
			close(stop) //主goroutine管理两个服务的生命周期
		}
	}

}

func StartHTTPServer(stop chan struct{}) error {
	server := http.Server{
		Addr:    ":6789",
		Handler: nil,
	}
	go func() {
		//控制StartHTTPServer的退出 stop由主goroutine进行操作
		<-stop
		server.Shutdown(context.Background())
	}()
	return server.ListenAndServe() //返回错误
}

//处理信号。。。这个信号没用过，不知道咋用，先暂时只处理这两个信号
//关于这个golang处理linux信号的处理，不知道老师能不能稍微提一下，以前没接触过
func HandleSingal(stop chan struct{}) error {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	go func() {
		//控制HandleSingal的退出
		<-stop
		c <- os.Interrupt
	}()

	<-c
	return errors.New("exit signal")
}
