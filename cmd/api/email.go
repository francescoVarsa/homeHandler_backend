package main

import (
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

func SendMessage(msg string, to string, smtpClient *mail.SMTPClient) {
	email := mail.NewMSG()
	email.SetFrom("admin@homeHandler.com").
		AddTo(to).
		SetSubject("New Go Email")

	email.SetBody(mail.TextHTML, msg)

	// always check error after send
	if email.Error != nil {
		log.Fatal(email.Error)
	}

	// Call Send and pass the client
	err := email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	}
}
