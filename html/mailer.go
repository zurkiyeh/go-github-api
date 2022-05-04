// Based off : https://medium.com/@dhanushgopinath/sending-html-emails-using-templates-in-golang-9e953ca32f3d

package html

import (
	"bytes"
	"html/template"
	"net/smtp"
)

//Request struct
type MailerRequest struct {
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *MailerRequest {
	return &MailerRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *MailerRequest) SendEmail(auth smtp.Auth, toEmailAddr, smtlHost, smtpPort string) (bool, error) {

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "zurkiyeh@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *MailerRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
