package mail

import (
	"crypto/tls"
	"net/mail"
	"net/smtp"
	"strconv"

	"bytes"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/emacsist/go-notify-center/code"
	"github.com/emacsist/go-notify-center/config"
	"github.com/emacsist/go-notify-center/message"
)

// Send : 发送 邮件
func Send(mail message.Email) (result message.CallbackData) {
	mail.From = config.Configuration.Email.UserName
	err := checkMail(mail)
	if err != nil {
		result = message.BuildCallbackData(mail.MessageID, code.NoToAddress, code.NoToAddressMsg+"\n"+err.Error(), mail.CallbackQueue, nil, mail.To, mail.From)
		return
	}

	c, err := smtpClient()
	if err != nil {
		result = message.BuildCallbackData(mail.MessageID, code.SMTPClientError, code.SMTPClientErrorInfo+"\n"+err.Error(), mail.CallbackQueue, nil, mail.To, mail.From)
		return
	}

	result.From = mail.From
	result.CallbackQueue = mail.CallbackQueue
	result.MessageID = mail.MessageID
	for _, toAddr := range mail.To {
		err = c.Reset()
		if err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}
		// To && From
		if err = c.Mail(config.Configuration.Email.UserName); err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}

		if err = c.Rcpt(toAddr); err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}

		// Data
		w, err := c.Data()
		if err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}
		_, err = w.Write(getMailMessage(toAddr, mail.Subject, mail.Body))
		if err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}
		err = w.Close()
		if err != nil {
			log.Errorf(err.Error())
			result.ToError = append(result.ToError, toAddr)
			result.ErrorCode = code.SMTPClientError
			result.ErrorInfo = toAddr + "==>" + code.SMTPClientErrorInfo + " ==> " + err.Error() + "\n"
			continue
		}
		result.ToOK = append(result.ToOK, toAddr)
	}
	if len(result.ErrorInfo) == 0 {
		result.ErrorCode = code.OK
		result.ErrorInfo = code.OKMsg
	}
	c.Quit()
	return
}

func checkMail(mail message.Email) error {
	if len(mail.To) == 0 {
		return fmt.Errorf("MessageID %s no to address", mail.MessageID)
	}
	return nil
}

// getMailMessage ： 构造邮件消息内容
func getMailMessage(toAddress string, subject string, body string) []byte {
	from := mail.Address{Address: config.Configuration.Email.UserName, Name: ""}
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
	host := config.Configuration.Email.Host + ":" + strconv.FormatInt(int64(config.Configuration.Email.Port), 10)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.Configuration.Email.Host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", host, tlsconfig)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	c, err := smtp.NewClient(conn, config.Configuration.Email.Host)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	auth := smtp.PlainAuth("", config.Configuration.Email.UserName, config.Configuration.Email.Password, config.Configuration.Email.Host)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return c, nil
}
