package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func login(cCtx *cli.Context) error {
	username := cCtx.String("username")

	fmt.Print("Enter your api key: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	user := user{
		Username: username,
		ApiKey:   string(b),
	}

	pl := newClient(user)

	spin := spinner.New(spinner.CharSets[2], 100*time.Millisecond)
	spin.Prefix = "please wait "
	spin.Suffix = " logging in"
	spin.Start()

	err = pl.authenticate()
	if err != nil {
		spin.Stop()
		return fmt.Errorf("error logging in: %s", err.Error())
	}

	spin.Stop()

	passphrase, err := readPassphrase(cCtx)
	if err != nil {
		return err
	}

	err = pl.writeAuth(passphrase)
	if err != nil {
		return err
	}

	return nil
}

func logout(cCtx *cli.Context) error {
	username := cCtx.String("username")

	passphrase, err := readPassphrase(cCtx)
	if err != nil {
		return err
	}

	_, err = readAuth(username, passphrase)
	if err != nil {
		return err
	}

	spin := spinner.New(spinner.CharSets[2], 100*time.Millisecond)
	spin.Prefix = "please wait "
	spin.Suffix = " logging out"
	spin.Start()

	if err != nil {
		return err
	}

	pldir, err := getPinlistDir()
	if err != nil {
		return err
	}

	exp := filepath.Join(pldir, username)
	err = os.Remove(exp)
	if err != nil {
		return err
	}

	return nil
}

func readAuth(username, passphrase string) (*client, error) {
	pldir, err := getPinlistDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(pldir, username)
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	plaintext, err := decrypt(passphrase, ciphertext)
	if err != nil {
		return nil, err
	}

	user := user{}
	err = json.Unmarshal([]byte(plaintext), &user)
	if err != nil {
		return nil,
			fmt.Errorf("error reading config: %s", err.Error())
	}

	return newClient(user), nil
}

func (c client) writeAuth(passphrase string) error {
	pldir, err := getPinlistDir()
	if err != nil {
		return err
	}

	b, err := json.Marshal(c.user)
	if err != nil {
		return err
	}

	path := filepath.Join(pldir, c.user.Username)
	ciphertext, err := encrypt(passphrase, b)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(ciphertext), 0700)
	if err != nil {
		return err
	}

	return nil
}

func readPassphrase(cCtx *cli.Context) (string, error) {
	passphrase := cCtx.String("passphrase")

	if strings.TrimSpace(passphrase) == "" {
		fmt.Print("Enter your passphrase: ")
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		passphrase = string(b)
	}

	if len(string(passphrase)) < 16 {
		return "", errors.New("passphrase must be 16 characters or more")
	}

	return passphrase, nil
}
