package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
	"time"
	"embed"
)

const (
	FromName            = "GopherSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitations.tmpl"
)

//go:embed "templates"
var FS embed.FS

type SMTPMailer struct {
	smtpHost  string
	smtpPort  string
	username  string
	password  string
}

func NewSMTPMailer(host, port, username, password string) *SMTPMailer {
	return &SMTPMailer{
		smtpHost:  host,
		smtpPort:  port,
		username:  username,
		password:  password,
	}
}

func (m *SMTPMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\n%s",
    m.username, email, subject.String(), body.String())

	auth := smtp.PlainAuth("", m.username, m.password, m.smtpHost)

	for i := 0; i < maxRetries; i++ {
		err = smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, m.username, []string{email}, []byte(message))
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d/%d: %v", email, i+1, maxRetries, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent successfully to %v", email)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
