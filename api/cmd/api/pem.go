package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func (app *application) generatePEM() error {

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

	return nil
}

func (app *application) encryptPEM(data []byte) ([]byte, error) {

	publicKeyPEM, err := os.ReadFile("./pem/public.pem")
	if err != nil {
		return nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (app *application) decryptPEM(data []byte) ([]byte, error) {

	privateKeyPEM, err := os.ReadFile("./pem/private.pem")
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)

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
