package main

import (
	"fmt"

	"github.com/emacsist/go-notify-center/mail"
	"github.com/emacsist/go-notify-center/message"
)

func main() {
	var m message.Email
	m.Body = "Hello World"
	m.CallbackQueue = "callback queue"
	m.MessageID = "1"
	m.Subject = "for test"
	m.To = append(m.To, "929168233@qq.com", "hellosdfsdfgm.ssdf.com")
	callBack := mail.Send(m)
	fmt.Printf("callback = %+v\n", callBack)
}
