package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
	"github.com/sirupsen/logrus"

	"github.com/medods-jwt-auth/config"
)

var (
	dialer *gomail.Dialer
)

func Init() {
	dialer = gomail.NewDialer(config.MAIL_HOST, int(config.MAIL_PORT), config.MAIL_USER, config.MAIL_PASSWORD)
}

func SendWarning(receiverEmail, ip string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.MAIL_USER)
	msg.SetHeader("To", receiverEmail)
	msg.SetHeader("Subject", "Medods account entered")
	msg.SetBody("text/html", fmt.Sprintf("Somebody login to your account from ip %s!", ip))

	if err := dialer.DialAndSend(msg); err != nil {
		logrus.Warning(err.Error())
		return err
	}

	return nil
}
