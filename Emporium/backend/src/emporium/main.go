package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	logfile, err := os.OpenFile("/var/log/emporium/emporium.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)
	http.HandleFunc("/", validate)
	http.ListenAndServe(":3000", nil)
	logfile.Close()
}

func validate(w http.ResponseWriter, r *http.Request) {

	var err error

	var payload []byte
	if r.Body != nil {
		payload, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	if err != nil {
		log.Printf("error reading payload %v", err)
		http.Error(w, "error rreading payload", http.StatusBadRequest)
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

	data, err := ioutil.ReadFile("/applepay/merchant.json")
	if err != nil {
		log.Printf("error preparing message %v", err)
		http.Error(w, "error preparing message", http.StatusInternalServerError)
		return
	}

	var msg bytes.Buffer
	err = json.Compact(&msg, data)
	if err != nil {
		log.Printf("error encoding message %v", err)
		http.Error(w, "error encoding message", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", params.URL, &msg)
	if err != nil {
		log.Printf("error preparing request %v", err)
		http.Error(w, "error preparing request", http.StatusInternalServerError)
		return
	}

	cert, err := tls.LoadX509KeyPair("/applepay/merchant.pem", "/applepay/merchant.pem")
	if err != nil {
		log.Printf("error preparing tls %v", err)
		http.Error(w, "error preparing tls", http.StatusInternalServerError)
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("error transporting msg %v", err)
		http.Error(w, "error transporting msg", http.StatusInternalServerError)
		return
	}

	log.Printf("handling res %v", res)

	// Defer closing of underlying connection so it can be re-used
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	var session []byte
	session, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error returning response %v", err)
		http.Error(w, "error returning response", http.StatusInternalServerError)
		return
	}

	log.Printf("returning session %v", session)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(session)))
	w.Write(session)

}
