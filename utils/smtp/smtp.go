package smtp

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/utils"
	gomail "gopkg.in/mail.v2"
)

func SendMail(to, from, fromPassword, subject, body string, config *config.Config) error {
	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", from)
	// Set E-Mail receivers
	m.SetHeader("To", to)
	// Set E-Mail subject
	m.SetHeader("Subject", subject)
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)
	// Settings for SMTP server
	emailPort, _ := strconv.ParseInt(config.EmailSmtpPort, 10, 0)
	d := gomail.NewDialer(config.EmailHostname, int(emailPort), from, fromPassword)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendSignInMail(to string, signInTime time.Time, config *config.Config, md *utils.ClientMetadata) error {
	msg := fmt.Sprintf("Se ha iniciado sesión en su cuenta el día %s.\n\nDispositivo: %s\nSistema Operativo: %s %s\n\nEnviamos este correo electrónico para asegurarnos de que haya sido usted.\nSi reconoce este inicio:\nNo es necesario hacer nada, puede ignorar este correo.\n\nGracias\nEl equipo de %s", signInTime.UTC(), *md.Model, *md.Platform, *md.SystemVersion, config.AppName)
	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", config.EmailAddress)
	// Set E-Mail receivers
	m.SetHeader("To", to)
	// Set E-Mail subject
	m.SetHeader("Subject", "Inicio de sesión")
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", msg)
	// Settings for SMTP server
	emailPort, _ := strconv.ParseInt(config.EmailSmtpPort, 10, 0)
	d := gomail.NewDialer(config.EmailHostname, int(emailPort), config.EmailAddress, config.EmailAddressPassword)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
