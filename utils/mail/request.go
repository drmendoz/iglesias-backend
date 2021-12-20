package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

type Request struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.Subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.Body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, usuario, r.To, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Body = buf.String()
	return nil
}
