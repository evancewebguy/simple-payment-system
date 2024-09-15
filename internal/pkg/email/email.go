package email

import (
	"crypto/tls"
	"mamlaka/config"

	"gopkg.in/gomail.v2"
)

type EmailMessage struct {
	To      string
	From    string
	Subject string
	Body    string
}

// SendEmail sends the email using an SMTP server
func SendEmail(emailMessage EmailMessage, conf config.EmailConfig) error {
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", conf.FromAddress)
	m.SetHeader("To", emailMessage.To)
	m.SetHeader("Subject", emailMessage.Subject)

	// Set email body (HTML or plain text)
	m.SetBody("text/html", emailMessage.Body)

	// Configure SMTP settings
	d := gomail.NewDialer(conf.SMTPServer, conf.SMTPPort, conf.Username, conf.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send email and handle error
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
