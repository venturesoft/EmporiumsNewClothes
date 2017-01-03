package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"miranda"
)

func main() {
	logfile, err := os.OpenFile("/var/log/emporium/emporium.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)
	http.HandleFunc("/getApplePaySession", validateMerchant)
	http.HandleFunc("/processPayment", processPayment)
	http.ListenAndServe(":3000", nil)
	logfile.Close()
}

func validateMerchant(w http.ResponseWriter, r *http.Request) {

	var err error

	var carmen miranda.MerchantValidationService
	carmen, err = miranda.CreateMerchantValidationService(30*time.Second, miranda.FileBasedMerchantValidationConfig{
		CertFilePath:        "/applepay/merchant.pem",
		RequestBodyFilePath: "/applepay/merchant.json",
	})
	if err != nil {
		log.Printf("error creating merchant validation service %v", err)
		http.Error(w, "invalid merchant", http.StatusInternalServerError)
		return
	}

	var payload []byte
	if r.Body != nil {
		payload, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	if err != nil {
		log.Printf("error reading payload %v", err)
		http.Error(w, "error reading payload", http.StatusBadRequest)
		return
	}

	var params = struct {
		URL string `json:"url"`
	}{}
	err = json.Unmarshal(payload, &params)
	if err == nil && params.URL == "" {
		err = errors.New("missing required parameter: url")
	}
	if err != nil {
		log.Printf("error parsing payload %v", err)
		http.Error(w, "error parsing payload", http.StatusBadRequest)
		return
	}

	var session json.RawMessage
	session, err = carmen.Dance(params.URL)
	if err != nil {
		log.Printf("error during merchant validation dance %v", err)
		http.Error(w, "merchant validation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(session)))
	w.Write(session)

}

func processPayment(w http.ResponseWriter, r *http.Request) {

	var err error

	var payload []byte
	if r.Body != nil {
		payload, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	if err != nil {
		log.Printf("error reading payload %v", err)
		http.Error(w, "error reading payload", http.StatusBadRequest)
		return
	}

	log.Printf("payload:\n %s", string(payload))

	// check if we are setup to test worldpay integration
	if f, err := os.Stat("/applepay/wap.json"); err == nil && !f.IsDir() {

		data, err := ioutil.ReadFile("/applepay/wap.json")
		if err != nil {
			log.Printf("error reading wap configuration %v", err)
			http.Error(w, "error reading wap configuration", http.StatusInternalServerError)
			return
		}

		var wapconfig = struct {
			MerchantCode string `json:"merchantCode"`
			Password     string `json:"password"`
		}{}
		err = json.Unmarshal(data, &wapconfig)
		if err == nil && (wapconfig.MerchantCode == "" || wapconfig.Password == "") {
			err = errors.New("missing required parameters: Merchant, Password")
		}
		if err != nil {
			log.Printf("error parsing wap configuration %v", err)
			http.Error(w, "error parsing wap configuration", http.StatusBadRequest)
			return
		}

		log.Printf("wap processing result :\n %s", wapProcess(wapconfig.MerchantCode, wapconfig.Password, payload))

	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "")

}
