package main

import (
	mail_checker "github.com/ngocchien/mail-checker"
	log "github.com/sirupsen/logrus"
)

func main() {
	emails := []string{
		"chiennn0104@hotmail.com",
		"chiennn0104123123123@hotmail.com",
	}
	checker := mail_checker.New(mail_checker.MailKindMicrosoft, mail_checker.Proxy{})
	for _, email := range emails {
		status := checker.Check(email)
		log.Infof("Email: %s, status: %+v", email, status)
	}
}
