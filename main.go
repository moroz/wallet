package main

import (
	"errors"
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"
)

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	pw, err := term.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

func getPassphrase() (string, error) {
	pw1, err := readPassword("Enter a password: ")
	if err != nil {
		return "", err
	}
	pw2, err := readPassword("Confirm password: ")
	if err != nil {
		return "", err
	}
	if pw1 != pw2 {
		return "", errors.New("password does not match confirmation")
	}
	return pw1, nil
}

func main() {
	passphrase, err := getPassphrase()
	if err != nil {
		log.Fatalf("Failed to create a passphrase for database: %s", err)
	}
	fmt.Println(string(passphrase))
}
