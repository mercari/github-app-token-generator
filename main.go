package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	b, ok := os.LookupEnv("APP_PRIVATE_KEY")
	if !ok {
		setFailed("Private key env not set", "Environment variable APP_PRIVATE_KEY is not set")
	}

	appID, ok := os.LookupEnv("APP_ID")
	if !ok {
		setFailed("App ID not set", "Environment variable APP_ID is not set")
	}

	appInstId, ok := os.LookupEnv("APP_INSTALLATION_ID")
	if !ok {
		setFailed("App installation ID not set", "Environment variable APP_INSTALLATION_ID is not set")
	}

	pemBytes, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		setFailed("Base64 decode failed", "PEM secret should be base64 encoded")
		return
	}

	key, err := loadPEMFromBytes(pemBytes)
	if err != nil {
		setFailed("Invalid PEM key", fmt.Sprintf("Unable to load PEM. err: %s", err))
	}

	jwt := issueJWTFromPEM(appID, key)
	token, err := getInstallationToken(appInstId, jwt)
	if err != nil {
		setFailed("Failed to get installation key", fmt.Sprintf("Unable to get intallation token. err: %s", err))
	}

	setOutput("token", token)
}

func loadPEMFromBytes(key []byte) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode(key)
	if b != nil {
		key = b.Bytes
	}

	parsedKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("private key should be a PKCS1 key; parse error: %v", err)
	}

	return parsedKey, nil
}

func setOutput(name, value string) {
	// print output to stdout
	fmt.Printf("::set-output name=%s::%s", name, value)
}

func setFailed(title, msg string) {
	// output error to stdout and exit 1
	fmt.Printf("::error title=%s::%s", title, msg)
	os.Exit(1)
}
