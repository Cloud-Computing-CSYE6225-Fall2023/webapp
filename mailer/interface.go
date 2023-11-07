package mailer

type Mailer interface {
	SendEmail(subject, body, recipient string)
}
