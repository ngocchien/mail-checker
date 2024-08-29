package main

import (
	"github.com/ngocchien/mail-checker"
	log "github.com/sirupsen/logrus"
)

func main() {
	emails := []string{
		"boy_codon_cangirl@yahoo.com",
		"boy_codon_cangirlxx1010100101k11k@yahoo.com",
	}
	checker := mail_checker.New(mail_checker.MailKindYahoo, mail_checker.Proxy{})
	for _, email := range emails {
		status := checker.Check(email)
		log.Infof("Email: %s, status: %+v", email, status)
	}
}
