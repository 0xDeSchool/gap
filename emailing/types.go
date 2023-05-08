package emailing

type EmailMessage struct {
	From       *string
	To         string
	Subject    string
	Body       string
	IsBodyHtml bool
}

type EmailOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailSender interface {
	Send(msg *EmailMessage) error
}
