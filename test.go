package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "NO-REPLY <no-reply@kora.moe>")
	m.SetHeader("To", "test <smtp+uvblyd3asit2npc5@mailtester.smtpserver.com>")
	m.SetHeader("Subject", "Account Action Required")
	m.SetBody("text/plain", "Your account is registered successfully, this is the code: 445283.")

	dialer := gomail.NewDialer("127.0.0.1", 25, "no-reply@kora.moe", "")
	err := dialer.DialAndSend(m)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
