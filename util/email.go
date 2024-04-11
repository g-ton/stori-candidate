package util

import (
	"fmt"
	"net/smtp"
)

func SendMail(to []string, msg []byte) error {
	// Sender data.
	from := "squallmagc@gmail.com"
	password := "ayafmzvpjgqybitg"
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email Sent Correctly!")
	return nil
}
