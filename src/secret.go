package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getAuth() (Auth, error) {
	err := godotenv.Load()
	if err != nil {
		return Auth{}, err
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	tgt := os.Getenv("TGT")
	if username == "" || password == "" {
		return Auth{}, fmt.Errorf("USERNAME or PASSWORD not found in .env file. Please add them.")
	}

	auth := Auth{
		Username: username,
		Password: password,
		TGT:      tgt,
	}
	return auth, err
}

// stupidest thing ever
// TODO stop doing this and store it properly
func setAuth(a Auth) error {
	err := godotenv.Write(map[string]string{
		"USERNAME": a.Username,
		"PASSWORD": a.Password,
		"TGT":      a.TGT,
	}, ".env")
	return err
}
