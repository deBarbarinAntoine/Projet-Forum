package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

func (api *API) encryptPEM(data []byte) ([]byte, error) {

	publicKeyBlock, _ := pem.Decode(api.pemKey)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		// DEBUG
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}

	hexEncrypted := hex.EncodeToString(ciphertext)

	return []byte(hexEncrypted), nil
}
