package main

import (
	"github.com/emacsist/go-notify-center/mail"
)

func main() {
	mail.Send([]string{"929168233@qq.com", "emacsist.yzy@gmail.com"}, "Hello", "Hello World")
}
