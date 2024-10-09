package email

import (
	"crypto/tls"
	"fmt"
	"server/src/helper"

	"gopkg.in/gomail.v2"
)

type Email struct {
	EMAIL    string
	PASSWORD string
}

func GeneratEmail() Email {
	e := helper.GetEnvVar("EMAIL")
	p := helper.GetEnvVar("EMAILPASSWD")

	return Email{
		EMAIL:    e,
		PASSWORD: p,
	}
}
func (e *Email) SendEmail(text string, sub string, email string) (helper.Status, error) {
	m := gomail.NewMessage()

	fmt.Println(email)
	// Set E-Mail sender
	m.SetHeader("From", e.EMAIL)

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", sub)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", "<h1> Luna </h1> "+text)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, e.EMAIL, e.PASSWORD)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	return helper.Status{STATUS: "Success"}, nil
}
