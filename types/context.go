package types

import (
	"github.com/shivasaicharanruthala/webapp/log"
	"github.com/shivasaicharanruthala/webapp/mailer"
	"github.com/shivasaicharanruthala/webapp/publish"
	"gopkg.in/alexcesaro/statsd.v2"
)

type Context struct {
	Logger        *log.CustomLogger
	Metrics       *statsd.Client
	MailerClient  mailer.Mailer
	PublishClient publish.SNSPublishAPI
}

func NewContext(logger *log.CustomLogger, metrics *statsd.Client, mailerClient mailer.Mailer, snsClient publish.SNSPublishAPI) *Context {
	return &Context{
		Logger:        logger,
		Metrics:       metrics,
		MailerClient:  mailerClient,
		PublishClient: snsClient,
	}
}
