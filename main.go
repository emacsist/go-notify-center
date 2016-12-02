package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"github.com/emacsist/go-notify-center/listener"
	"github.com/emacsist/go-notify-center/message"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	var m message.Email
	m.Body = "hello"
	m.CallbackQueue = "icallback"
	m.MessageID = "myid"
	m.Subject = "subject"
	m.To = append(m.To, "929168233@qq.com", "sldfkjsdfsl@sdfsf.c")

	b, _ := json.Marshal(m)
	fmt.Printf("%v\n", string(b))

	// sigs 表示将收到的信号放到这个管道中。后面的参数表示你想处理的系统信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//开启一个Go routine来监听信号
	go func() {
		sig := <-sigs
		fmt.Printf("收到信号: %v\n", sig)
		done <- true
	}()

	//这里添加你的程序的功能（监听器，处理器等）
	fmt.Println("程序启动完毕，等待接收信号中...")
	<-done
	exit()
	fmt.Println("成功退出.")
}

func exit() {
	listener.Close()
	fmt.Printf("执行清理OK\n")
}
