package mail

import (
	"fmt"
	"net/smtp"

	"github.com/g-ton/stori-candidate/util"
)

type Mail interface {
	// The interface defining the abstract methods
	// Maybe It would be better to set this interface in another file, however, it was put here to keep the things simple
	SendMail(to []string, msg []byte) error
}

// - Service
type MailService struct {
	Credentials mailCredentials
}

func NewMail(c util.Config) *MailService {
	return &MailService{
		Credentials: setCredentials(c),
	}
}

type mailCredentials struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func setCredentials(c util.Config) mailCredentials {
	return mailCredentials{
		from:     c.EmailFrom,
		password: c.EmailPassword,
		smtpHost: c.EmailSmtpHost,
		smtpPort: c.EmailSmtpPort,
	}
}

// - Implementation of methods
func (es *MailService) SendMail(to []string, msg []byte) error {
	// Authentication.
	auth := smtp.PlainAuth("", es.Credentials.from, es.Credentials.password, es.Credentials.smtpHost)

	// Sending email.
	err := smtp.SendMail(es.Credentials.smtpHost+":"+es.Credentials.smtpPort, auth, es.Credentials.from, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email Sent Correctly!")
	return nil
}
