package mailer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type client struct {
	sender *mailgun.MailgunImpl
}

func New() Mailer {
	// Create an instance of the Mailgun Client
	domain := os.Getenv("MAILGUN_API_KEY")
	apiKey := os.Getenv("MAILGUN_DOMAIN")
	ms := mailgun.NewMailgun(domain, apiKey)

	return &client{
		sender: ms,
	}
}

func (mc *client) SendEmail(subject, body, recipient string) {
	sender := os.Getenv("MAILGUN_SENDER_EMAIL")

	// The message object allows you to add attachments and Bcc recipients
	message := mc.sender.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mc.sender.Send(ctx, message)
	if err != nil {
		fmt.Printf("ERROR: Sending email with error %v", err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
