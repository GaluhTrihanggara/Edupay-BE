package mail

import (
	"Edupay/config"
	"Edupay/model"
	"bytes"
	"fmt"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type MailUseCase interface {
	SendingMail(payload model.PayloadMail)
}

const (
	IMAGE_LOGO     = "https://res.cloudinary.com/duoehn6px/image/upload/v1687359124/ppob/jzvolgffbkgtspzqcs7r.png"
	TEXT_LOGO      = "https://res.cloudinary.com/duoehn6px/image/upload/v1687359227/ppob/jxzbamfzcap7hy11yvgb.png"
	TWITTER_LOGO   = "https://res.cloudinary.com/ddleabcu2/image/upload/v1687440327/ppob/tw5v6ducw6ee76ghka9c.png"
	INSTAGRAM_LOGO = "https://res.cloudinary.com/ddleabcu2/image/upload/v1687440400/ppob/fd8qkyywqr6broxxzt4i.png"
	FACEBOOK_LOGO  = "https://res.cloudinary.com/ddleabcu2/image/upload/v1687440451/ppob/nuocmo2f6qejfruiedqe.png"
)

func SendingMail(payload model.PayloadMail) {
	tmpl, err := template.ParseFiles("usecase/mail/invoice.html")
	if err != nil {
		log.Fatal(err.Error())
	}

	var emailBody bytes.Buffer
	err = tmpl.Execute(&emailBody, payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.AppConfig.SenderEmail)
	mailer.SetHeader("To", payload.RecipentEmail)
	mailer.SetHeader("Subject", fmt.Sprintf("Invoice Pembayaran %s", payload.ProductType))
	mailer.SetBody("text/html", emailBody.String())

	dialer := gomail.NewDialer(
		config.AppConfig.SmtpHost,
		config.AppConfig.SmtpPort,
		config.AppConfig.SenderEmail,
		config.AppConfig.EmailPassword,
	)
	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
