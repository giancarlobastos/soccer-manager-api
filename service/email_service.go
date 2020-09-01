package service

import (
	"fmt"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"net/smtp"
)

type EmailService struct {
	auth *smtp.Auth
}

func NewEmailService() *EmailService {
	auth := smtp.PlainAuth("", "soccer.manager.api@gmail.com", "soccer.manager", "smtp.gmail.com")
	return &EmailService{
		auth: &auth,
	}
}

func (es *EmailService) SendVerificationEmail(account *domain.Account) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Verify your account!\n"
	body := fmt.Sprintf("<html><body><p>Hi %s!</p><a href='http://localhost:8080/verify-account?token=%s'>Click here to verify your account!</a></body></html>",
		account.FirstName, account.VerificationToken)
	msg := []byte(subject + mime + body)
	addr := "smtp.gmail.com:587"

	return smtp.SendMail(addr, *es.auth, "soccer.manager.api@gmail.com", []string{account.Username}, msg)
}
