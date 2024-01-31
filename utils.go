package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func getPinlistDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pldir := filepath.Join(homedir, ".pinlist")
	if _, err := os.Stat(pldir); os.IsNotExist(err) {
		err := os.MkdirAll(pldir, 0700)
		if err != nil {
			return "", err
		}
	}

	return pldir, nil
}

func encrypt(key string, plaintext []byte) (string, error) {
	sha := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sha[0:])
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key string, ciphertext []byte) (string, error) {
	c, err := base64.URLEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	sha := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sha[0:])
	if err != nil {
		return "", nil
	}

	if len(c) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := c[:aes.BlockSize]
	c = c[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(c, c)

	return string(c), nil
}
