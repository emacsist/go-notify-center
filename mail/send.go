package mail

import (
	"crypto/tls"
	"net/mail"
	"net/smtp"
	"strconv"

	"bytes"

	log "github.com/Sirupsen/logrus"
)

// Send : 发送 邮件
func Send(toAddress []string, subject string, body string) error {

	c, err := smtpClient()

	for _, toAddr := range toAddress {
		// To && From
		if err = c.Mail(MailConfig.UserName); err != nil {
			log.Panic(err)
			return err
		}

		if err = c.Rcpt(toAddr); err != nil {
			log.Panic(err)
			return err
		}

		// Data
		w, err := c.Data()
		if err != nil {
			log.Panic(err)
			return err
		}
		_, err = w.Write(getMailMessage(toAddr, subject, body))
		if err != nil {
			log.Panic(err)
			return err
		}
		err = w.Close()
		if err != nil {
			log.Panic(err)
			return err
		}
	}
	c.Quit()
	return nil
}

// getMailMessage ： 构造邮件消息内容
func getMailMessage(toAddress string, subject string, body string) []byte {
	from := mail.Address{Address: MailConfig.UserName, Name: ""}
	to := mail.Address{Address: toAddress, Name: toAddress}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	var message bytes.Buffer
	// Setup message
	for k, v := range headers {
		message.WriteString(k)
		message.WriteString(":")
		message.WriteString(v)
		message.WriteString("\r\n")
	}
	message.WriteString("\r\n")
	message.WriteString(body)
	return message.Bytes()
}

// smtpClient : 相当于一个客户端（已经认证OK的，如果成功的话）
func smtpClient() (*smtp.Client, error) {
	host := MailConfig.Host + ":" + strconv.FormatInt(int64(MailConfig.Port), 10)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         MailConfig.Host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", host, tlsconfig)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	c, err := smtp.NewClient(conn, MailConfig.Host)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	auth := smtp.PlainAuth("", MailConfig.UserName, MailConfig.Password, MailConfig.Host)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
		return nil, err
	}
	return c, nil
}
