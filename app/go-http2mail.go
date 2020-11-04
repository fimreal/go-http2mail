package main

import (
	config "go-http2mail/pkg/config"
	serve "go-http2mail/pkg/http"
	mail "go-http2mail/pkg/mail"

	"github.com/spf13/viper"
)

func main() {
	//
	config.LoadConfigs()
	serve.HandleRequests(mail.SendMail, ":"+viper.GetString("port"), "/mailme")
}
