package listener

import (
	"encoding/json"
	"time"

	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/emacsist/go-notify-center/config"
	"github.com/emacsist/go-notify-center/mail"
	"github.com/emacsist/go-notify-center/message"
	"github.com/streadway/amqp"
)

var (
	// conn : rabbitmq ：连接对象
	conn             *amqp.Connection
	rabbitCloseError chan *amqp.Error
)

// Close :关闭
func Close() {
	if config.Configuration.Callback.IsOK {
		if conn != nil && !conn.IsClosed() {
			err := conn.Close()
			if err != nil {
				log.Errorf("close rabbit connection error %v", err.Error())
			} else {
				log.Infof("rabbit connection is closed.")
			}
		} else {
			log.Infof("rabbit connection is closed.")
		}
	}
}

func init() {

	if config.Configuration.Callback.IsOK {
		go func() {
			rabbitCloseError = make(chan *amqp.Error)
			rabbitCloseError <- amqp.ErrClosed
		}()
		go rabbitConnector(config.Configuration.Callback.AMQPURL)
	}
}

func connectToRabbitMQ(uri string) *amqp.Connection {
	for {
		conn, err := amqp.Dial(uri)
		if err == nil {
			return conn
		}
		log.Error(err.Error())
		log.Warnf("Trying to reconnect to RabbitMQ at %v", uri)
		time.Sleep(500 * time.Millisecond)
	}
}
func rabbitConnector(uri string) {
	for {
		rabbitErr := <-rabbitCloseError
		if rabbitErr != nil {
			log.Errorf("Connecting to %v\n", uri)
			conn = connectToRabbitMQ(uri)
			conn.NotifyClose(rabbitCloseError)
			listen()
		}
	}
}

func listen() {
	conn = connectToRabbitMQ(config.Configuration.Callback.AMQPURL)
	ch, err := conn.Channel()
	if err != nil {
		panic("get connection channel error " + err.Error())
	}
	defer ch.Close()
	defer conn.Close()

	q, err := ch.QueueDeclare(
		config.Configuration.Callback.QueueName, // name
		true,  // durablelogrus
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic("fail to declar a queue : " + config.Configuration.Callback.QueueName + ", why = " + err.Error())
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"go-notify-center-mail-consumer", // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		panic("Failed to register a consumer " + err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(config.Configuration.Callback.Worker)
	for i := 0; i < config.Configuration.Callback.Worker; i++ {
		go func() {
			routineMqChannel, err := conn.Channel()
			if err != nil {
				panic("get connection channel error " + err.Error())
			}
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				email := message.Email{}
				err := json.Unmarshal(d.Body, &email)
				if err != nil {
					log.Errorf("parse %v to json error . why = %v", string(d.Body), err.Error())
					continue
				}
				callback := mail.Send(email)
				pushMessage(routineMqChannel, callback)
			}
			log.Infof("exit rabbit worker ...")
			wg.Done()
			routineMqChannel.Close()
		}()
	}
	log.Printf(" [*] mail listener start %v worker done. daemon...", config.Configuration.Callback.Worker)
	wg.Wait()
	go func() {
		rabbitCloseError = make(chan *amqp.Error)
		rabbitCloseError <- amqp.ErrClosed
	}()
}

func pushMessage(ch *amqp.Channel, callback message.CallbackData) {

	body, err := json.Marshal(callback)
	if err != nil {
		log.Errorf("parse %+v to json error, why = %v", callback, err.Error())
		return
	}

	q, err := ch.QueueDeclare(callback.CallbackQueue, // name
		true,  // durablelogrus
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Errorf("declare queue [%v] error, why = %v", callback.CallbackQueue, err.Error())
		return
	}

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "utf-8",
		Body:            body,
		DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
		Priority:        0,               // 0-9
		// a bunch of application/implementation-specific fields
	})
	if err != nil {
		log.Errorf("push callback message [%+v] to queue [%v] error, why = [%v]", callback, callback.CallbackQueue, err.Error())
	} else {
		log.Infof("push callback message [%+v] to queue [%v] OK", callback, q.Name)
	}
}
