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
func Send(toAddress string, subject string, body string) error {
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

	// Connect to the SMTP Server
	servername := MailConfig.Host + ":" + strconv.FormatInt(int64(MailConfig.Port), 10)

	auth := smtp.PlainAuth("", MailConfig.UserName, MailConfig.Password, MailConfig.Host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         MailConfig.Host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
		return err
	}

	c, err := smtp.NewClient(conn, MailConfig.Host)
	if err != nil {
		log.Panic(err)
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
		return err
	}

	_, err = w.Write(message.Bytes())
	if err != nil {
		log.Panic(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
		return err
	}

	c.Quit()
	return nil
}
