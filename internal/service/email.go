package service

import (
	"e-wallet/domain"
	"e-wallet/internal/config"
	"net/smtp"
)

type emailService struct {
	cnf *config.Config
}

func NewEmail(cnf *config.Config) domain.EmailService {
	return &emailService{
		cnf: cnf,
	}
}

func (e emailService) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.cnf.Email.User, e.cnf.Email.Password, e.cnf.Email.Host)
	msg := []byte(
		"From: Pt. Pencari Cinta Sejati <" + e.cnf.Email.User + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body,
	)
	return smtp.SendMail(e.cnf.Email.Host+":"+e.cnf.Email.Port, auth, e.cnf.Email.User, []string{to}, msg)
}
