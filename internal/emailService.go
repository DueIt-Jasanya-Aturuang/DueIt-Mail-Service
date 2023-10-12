package internal

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"

	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/config"
	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/template"
)

type EmailService interface {
	SendSmtp(payload []byte) error
	SendGOMAIL(payload []byte) error
}

type EmailServiceImpl struct {
	Template *template.EmailTemplateImpl
}

func NewEmailServiceImpl(template *template.EmailTemplateImpl) *EmailServiceImpl {
	return &EmailServiceImpl{
		Template: template,
	}
}

func (e *EmailServiceImpl) SendSmtp(payload []byte) error {
	var mail Email
	if err := json.Unmarshal(payload, &mail); err != nil {
		return err
	}

	to := []string{mail["to"]}
	cc := []string{"jasanya.tech@gmail.com"}

	templateMailValue := []byte(mail["value"])

	smtpAuth := smtp.PlainAuth("jasanya auth", config.Get().Mail.Address, config.Get().Mail.Pass, config.Get().Mail.Host)
	smtpAddrs := fmt.Sprintf("%s:%d", config.Get().Mail.Host, config.Get().Mail.Port)

	if err := smtp.SendMail(smtpAddrs, smtpAuth, config.Get().Mail.Address, append(to, cc...), templateMailValue); err != nil {
		return err
	}

	return nil
}

func (e *EmailServiceImpl) SendGOMAIL(payload []byte) error {
	var mail Email
	if err := json.Unmarshal(payload, &mail); err != nil {
		return err
	}
	log.Info().Msgf("%v", mail)

	var templateMailValue string
	if mail["type"] == "activasi-account" {
		templateBuffer := e.Template.CodeOTP(mail)
		templateMailValue = templateBuffer.String()
	} else if mail["type"] == "forgot-password" {
		templateBuffer := e.Template.ForgotPassword(mail)
		templateMailValue = templateBuffer.String()
	} else {
		templateMailValue = mail["message"]
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.Get().Mail.Sender)
	mailer.SetHeader("To", mail["to"], "jasanya@gmail.com")
	mailer.SetAddressHeader("Cc", "jasanya.tech@gmail.com", "jasanyatech")
	mailer.SetHeader("Subject", mail["type"])
	mailer.SetBody("text/html", templateMailValue)
	// mailer.AddAlternative("text/plain", html2text.HTML2Text(template))

	dialer := gomail.NewDialer(
		config.Get().Mail.Host,
		config.Get().Mail.Port,
		config.Get().Mail.Address,
		config.Get().Mail.Pass,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	log.Info().Msg("mail success send")

	return nil
}
