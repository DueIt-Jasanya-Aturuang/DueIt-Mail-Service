package template

import (
	"bytes"

	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/modules/entities"
	"github.com/rs/zerolog/log"
)

type EmailTemplateImpl struct{}

func NewEmailTemplateImpl() *EmailTemplateImpl {
	return &EmailTemplateImpl{}
}

func (t *EmailTemplateImpl) CodeOTP(data *entities.Email) bytes.Buffer {
	var body bytes.Buffer

	template, err := ParseTemplateDir("internal/template/html")
	if err != nil {
		log.Err(err).Msg("could not parse template")
	}

	template.ExecuteTemplate(&body, "codeOtp.html", &data)
	return body
}

func (t *EmailTemplateImpl) ForgotPassword(data *entities.Email) bytes.Buffer {
	var body bytes.Buffer

	template, err := ParseTemplateDir("internal/template/html")
	if err != nil {
		log.Err(err).Msg("could not parse template")
	}

	template.ExecuteTemplate(&body, "forgotPassword.html", &data)
	return body
}
