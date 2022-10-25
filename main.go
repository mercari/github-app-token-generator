package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
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

	githubOutputFile, err := openGitHubOutput(os.Getenv("GITHUB_OUTPUT"))
	if err != nil {
		setFailed("Failed to set the output 'token'", fmt.Sprintf("Unable to open GITHUB_OUTPUT. err: %s", err))
		return
	}
	defer githubOutputFile.Close()

	if err := setOutput(githubOutputFile, "token", token); err != nil {
		setFailed("Failed to set the output 'token'", fmt.Sprintf("Unable to write the token to GITHUB_OUTPUT. err: %s", err))
	}
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

func openGitHubOutput(p string) (io.WriteCloser, error) {
	return os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

func setOutput(file io.Writer, name, value string) error {
	if _, err := file.Write([]byte(fmt.Sprintf("%s=%s\n", name, value))); err != nil {
		return err
	}
	return nil
}

func setFailed(title, msg string) {
	// output error to stdout and exit 1
	fmt.Printf("::error title=%s::%s", title, msg)
	os.Exit(1)
}
