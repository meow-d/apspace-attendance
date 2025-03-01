package main

import (
	"github.com/zalando/go-keyring"
)

const Service string = "apspace"

func getAuth() Auth {
	username, _ := keyring.Get(Service, "username")
	password, _ := keyring.Get(Service, "password")
	tgt, _ := keyring.Get(Service, "tgt")

	auth := Auth{
		Username: username,
		Password: password,
		TGT:      tgt,
	}
	return auth
}

func setAuth(a Auth) error {
	err := keyring.Set(Service, "username", a.Username)
	if err != nil {
		return err
	}

	err = keyring.Set(Service, "password", a.Password)
	if err != nil {
		return err
	}

	err = keyring.Set(Service, "tgt", a.TGT)
	if err != nil {
		return err
	}

	return nil
}
