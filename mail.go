// Package ini provides structures and methods for parsing and managing INI files with special enhancements.
package mail

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type EmailSenderOptionSetter func(*EmailSender)

func WithHost(host string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Host = host
	}
}

func WithPort(port int) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Port = port
	}
}

func WithUsername(username string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Username = username
	}
}

func WithPassword(password string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Password = password
	}
}

func WithFrom(from string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.From = from
	}
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Email struct {
	To      []string
	Subject string
	Body    string
}

type EmailSender struct {
	config EmailConfig
}

func NewEmailSender(options ...EmailSenderOptionSetter) *EmailSender {
	es := &EmailSender{}

	for _, option := range options {
		option(es)
	}

	return es
}

func (e *EmailSender) Send(email Email) error {
	if len(email.To) == 0 {
		return errors.New("no recipients specified")
	}

	m := gomail.NewMessage()

	m.SetHeader("From", e.config.From)
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Body)

	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
	return d.DialAndSend(m)
}
