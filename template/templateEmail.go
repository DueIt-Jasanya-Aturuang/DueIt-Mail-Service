package template

import (
	"bytes"

	"github.com/rs/zerolog/log"
)

type EmailTemplateImpl struct{}

func NewEmailTemplateImpl() *EmailTemplateImpl {
	return &EmailTemplateImpl{}
}

func (t *EmailTemplateImpl) CodeOTP(data map[string]string) bytes.Buffer {
	var body bytes.Buffer

	template, err := ParseTemplateDir("template/html")
	if err != nil {
		log.Err(err).Msg("could not parse template")
	}

	err = template.ExecuteTemplate(&body, "codeOtp.html", &data)
	if err != nil {
		log.Err(err).Msg("failed to execute template")
	}
	return body
}

func (t *EmailTemplateImpl) ForgotPassword(data map[string]string) bytes.Buffer {
	var body bytes.Buffer

	template, err := ParseTemplateDir("template/html")
	if err != nil {
		log.Err(err).Msg("could not parse template")
	}

	err = template.ExecuteTemplate(&body, "forgotPassword.html", &data)
	if err != nil {
		if err != nil {
			log.Err(err).Msg("failed to execute template")
		}
	}
	return body
}
