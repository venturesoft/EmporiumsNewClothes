// Package miranda is a server side dance partner for Apple Pay Merchant Validation.
//
// It is named after Carmen Miranda, the dancer in the tutti-frutti hat.
package miranda

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// MerchantValidationConfig defines the configuration required
type MerchantValidationConfig interface {
	// LoadCert returns the merchant identity certificate used for the two-way TLS, see https://developer.apple.com/reference/applepayjs/applepaysession#2168856
	LoadCert() (tls.Certificate, error)
	// LoadRequestBody returns the JSON encoded merchant request payload, see https://developer.apple.com/reference/applepayjs/applepaysession#2168856
	LoadRequestBody() (json.RawMessage, error)
}

// FileBasedMerchantValidationConfig provides a simple file based implementation of the required configuration
type FileBasedMerchantValidationConfig struct {
	// CertFilePath defines the location of a PEM encoded certificate file (the private key must be included in the file)
	CertFilePath string
	// RequestBodyFilePath specifies the location of a JSON encoded file defining the merchant request payload
	RequestBodyFilePath string
}

// MerchantValidationService defines the server side steps in the Apple Pay Merchant Validation dance.
//
// A partner is required for this dance in order to perform the client side steps and provide the validation URL.
//
// See https://developer.apple.com/reference/applepayjs/applepaysession#2166532
type MerchantValidationService interface {
	// Dance performs the server side steps in the Apple Pay Merchant Validation dance returning a merchant session object.
	Dance(url string) (json.RawMessage, error)
}

type merchantValidation struct {
	reqBody *bytes.Buffer
	client  *http.Client
}

// LoadCert reads and parses a public/private key pair from the CertFilePath. The file must contain PEM encoded data.
func (c FileBasedMerchantValidationConfig) LoadCert() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.CertFilePath, c.CertFilePath)
}

// LoadRequestBody reads the file named by RequestBodyFilePath and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func (c FileBasedMerchantValidationConfig) LoadRequestBody() (json.RawMessage, error) {
	return ioutil.ReadFile(c.RequestBodyFilePath)
}

// CreateMerchantValidationService initiates the MerchantValidationService
func CreateMerchantValidationService(timeout time.Duration, conf MerchantValidationConfig) (MerchantValidationService, error) {

	data, err := conf.LoadRequestBody()
	if err != nil {
		return nil, fmt.Errorf("error loading request body %v", err)
	}

	var body bytes.Buffer
	err = json.Compact(&body, data)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body %v", err)
	}

	cert, err := conf.LoadCert()
	if err != nil {
		return nil, fmt.Errorf("error preparing tls %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}
	return merchantValidation{&body, client}, nil
}

func (s merchantValidation) Dance(url string) (json.RawMessage, error) {

	req, err := http.NewRequest("POST", url, s.reqBody)
	if err != nil {
		return nil, fmt.Errorf("error preparing request %v", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error transporting request %v", err)
	}

	// Defer closing of underlying connection so it can be re-used
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status transporting request %s", res.Status)
	}

	var session []byte
	session, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error returning response %v", err)
	}

	return session, nil
}
