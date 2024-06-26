package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (app *application) getPEM() (err error) {

	if fileExists("./pem/private.pem") && fileExists("./pem/public.pem") {
		app.config.pem.publicKey, err = os.ReadFile("./pem/public.pem")
		if err != nil {
			return err
		}
		app.config.pem.privateKey, err = os.ReadFile("./pem/private.pem")
		if err != nil {
			return err
		}
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile("./pem/private.pem", privateKeyPEM, 0644)
	if err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile("./pem/public.pem", publicKeyPEM, 0644)
	if err != nil {
		return err
	}

	app.config.pem.privateKey = privateKeyPEM
	app.config.pem.publicKey = publicKeyPEM

	return nil
}

func (app *application) decryptPEM(data []byte) ([]byte, error) {

	privateKeyBlock, _ := pem.Decode(app.config.pem.privateKey)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
