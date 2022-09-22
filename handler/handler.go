package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"gin-test/util/flag"
	"gin-test/util/log"
	"io/ioutil"
)

const (
	Message = "message"
)

func Init() error {
	err := createKeys()
	if err != nil {
		log.Logger("create key failed: %s\n", err)
		return err
	}
	return nil
}

func createKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	err = ioutil.WriteFile(flag.PriPath, pem.EncodeToMemory(priBlock), 0644)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err = ioutil.WriteFile(flag.PubPath, pem.EncodeToMemory(publicBlock), 0644)
	if err != nil {
		return err
	}
	flag.PubKey = string(pem.EncodeToMemory(publicBlock))
	flag.PriKey = privateKey
	return nil
}
