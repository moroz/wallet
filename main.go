package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"syscall"

	"github.com/tyler-smith/go-bip39"
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

func createDB(bin, dbPath, passphrase string) error {
	stdin := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", passphrase, passphrase))

	var stdout bytes.Buffer
	cmd := exec.Command(bin, "db-create", dbPath, "-p")
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func addEntryToDB(bin, dbPath, passphrase, entryName, password string) error {
	stdin := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", passphrase, password))

	var stdout bytes.Buffer
	cmd := exec.Command(bin, "add", dbPath, "-p", entryName)
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func main() {
	bin, err := exec.LookPath("keepassxc-cli")
	if err != nil {
		log.Fatal(err)
	}

	passphrase, err := getPassphrase()
	if err != nil {
		log.Fatalf("Failed to create a passphrase for database: %s", err)
	}

	err = createDB(bin, "./wallet.kdbx", passphrase)

	var seed = make([]byte, 32)
	_, err = rand.Read(seed)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, _ := bip39.NewMnemonic(seed)

	err = addEntryToDB(bin, "./wallet.kdbx", passphrase, "WALLET_SEED", base64.StdEncoding.EncodeToString(seed))
	if err != nil {
		log.Fatal(err)
	}

	err = addEntryToDB(bin, "./wallet.kdbx", passphrase, "WALLET_MNEMONIC", mnemonic)
	if err != nil {
		log.Fatal(err)
	}
}
