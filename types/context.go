package types

import (
	"github.com/shivasaicharanruthala/webapp/log"
	"github.com/shivasaicharanruthala/webapp/mailer"
	"gopkg.in/alexcesaro/statsd.v2"
)

type Context struct {
	Logger       *log.CustomLogger
	Metrics      *statsd.Client
	MailerClient mailer.Mailer
}

func NewContext(logger *log.CustomLogger, metrics *statsd.Client, mailerClient mailer.Mailer) *Context {
	return &Context{
		Logger:       logger,
		Metrics:      metrics,
		MailerClient: mailerClient,
	}
}
