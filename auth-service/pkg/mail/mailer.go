package mail

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendWelcomeEmail(toEmail, toName string) error
}

type smtpMailer struct {
	dialer *gomail.Dialer
	sender string
}

func NewSmtpMailer(host string, port string, user, password, sender string) Mailer {
	p, _ := strconv.Atoi(port)
	dialer := gomail.NewDialer(host, p, user, password)
	return &smtpMailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m *smtpMailer) SendWelcomeEmail(toEmail, toName string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Selamat Datang di Aplikasi Kami!")
	msg.SetBody("text/html", "Halo <b>"+toName+"</b>,<br><br>Terima kasih telah mendaftar.")

	// Kirim email
	if err := m.dialer.DialAndSend(msg); err != nil {
		log.Printf("Failed to send welcome email to %s: %v", toEmail, err)
	} else {
		log.Printf("Welcome email sent to %s", toEmail)
	}
	
	return nil
}
