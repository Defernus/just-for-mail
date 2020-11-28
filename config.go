package main

import (
	"log"
	"os"
)

const (
	envIsAutorizationRequired = "IS_AUTORIZATION_REQUIRED"
	envMailUser = "MAIL_USER"
	envMailPassword = "MAIL_PASSWORD"
)


type configData struct {
	isAutorisationRequired bool
	user                   string
	password               string
}

var config = configData {
	isAutorisationRequired: true,
	user: "",
	password: "",
}

func init() {
	isAutorisationRequired := os.Getenv(envIsAutorizationRequired)
	user := os.Getenv(envMailUser)
	password := os.Getenv(envMailPassword)

	if(isAutorisationRequired == "0") {
		log.Println("autorization not required")
		config.isAutorisationRequired = false
	}

	log.Println("autorization required")

	// if len(user) == 0 {
	// 	log.Fatalf("required %s", envMailUser)
	// }
	// if len(password) == 0 {
	// 	log.Fatalf("required %s", envMailPassword)
	// }

	config.user = user
	config.password = password
}
