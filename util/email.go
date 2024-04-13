package util

import (
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
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

/*
	GetAbsRootPath gets the absolute path for the project, doesn't matter where the file is executed

always it will returns the absolute project path, i.e :
if the file is in "stori-candidate/api/helper" the absolute path will be "/home/user/stori-candidate"
if the file is in "stori-candidate/api" the absolute path will be "/home/user/stori-candidate"
Please if you think about change the name of the project, you need to change it inside this func too
*/
func GetAbsRootPath() string {
	pwd, _ := os.Getwd()
	folder := strings.Split(pwd, string(filepath.Separator))
	newPath := ""
	for _, f := range folder {
		if f == "" {
			continue
		}
		newPath += string(filepath.Separator) + f
		if f == "stori-candidate" {
			break
		}
	}
	return newPath
}
