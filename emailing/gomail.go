package emailing

import (
	"gopkg.in/gomail.v2"
)

type GoMailSender struct {
	dialer *gomail.Dialer
	opts   *EmailOptions
}

func NewGoMailSender(opts *EmailOptions) *GoMailSender {
	d := gomail.NewDialer(opts.Host, opts.Port, opts.Username, opts.Password)
	d.Auth = LoginAuth(opts.Username, opts.Password)
	sender := GoMailSender{
		dialer: d,
		opts:   opts,
	}
	return &sender
}

func (gs *GoMailSender) Send(msg *EmailMessage) error {
	m := gomail.NewMessage()
	if msg.From == nil {
		m.SetHeader("From", gs.opts.From)
	} else {
		m.SetHeader("From", *msg.From)
	}
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	if msg.IsBodyHtml {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}
	return gs.dialer.DialAndSend(m)
}
