package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type wapTransaction struct {
	Payment             string // base 64 encoded JSON
	PaymentObject       PaymentObject
	OrderCode           string
	OrderDescription    string
	ShopperLanguageCode string
	AmountValue         string
	AmountCurrencyCode  string
	AmountExponent      string
	MerchantCode        string
}

type PaymentObject struct {
	Token           PaymentToken   `json:"token"`
	BillingContact  PaymentContact `json:"billingContact"`
	ShippingContact PaymentContact `json:"shippingContact"`
}

type PaymentToken struct {
	PaymentMethod struct {
		Network     string `json:"network"`
		Type        string `json:"type"`
		DisplayName string `json:"displayName"`
	} `json:"paymentMethod"`
	TransactionIdentifier string `json:"transactionIdentifier"`
	PaymentData           struct {
		Data      string `json:"data"`
		Signature string `json:"signature"`
		Header    struct {
			PublicKeyHash      string `json:"publicKeyHash"`
			EphemeralPublicKey string `json:"ephemeralPublicKey"`
			TransactionID      string `json:"transactionId"`
		} `json:"header"`
		Version string `json:"version"`
	} `json:"paymentData"`
}

type PaymentContact struct {
	EmailAddress       string   `json:"emailAddress"`
	PhoneNumber        string   `json:"phoneNumber"`
	FamilyName         string   `json:"familyName"`
	GivenName          string   `json:"givenName"`
	AddressLines       []string `json:"addressLines"`
	Locality           string   `json:"locality"`
	PostalCode         string   `json:"postalCode"`
	AdministrativeArea string   `json:"administrativeArea"`
	Country            string   `json:"country"`
	CountryCode        string   `json:"countryCode"`
}

var wapTemplate *template.Template
var wapTemplateErr error

func init() {
	wapTemplate, wapTemplateErr = template.New("wapTemplate").Parse(`<?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE paymentService PUBLIC "-//WorldPay/DTD WorldPay PaymentService v1//EN" "http://dtd.worldpay.com/paymentService_v1.dtd">
    <paymentService version="1.4" merchantCode="{{ .MerchantCode }}">
      <submit>
        <order orderCode="{{ .OrderCode }}" shopperLanguageCode="{{ .ShopperLanguageCode }}">
          <description>{{ .OrderDescription }}</description>
          <amount value="{{ .AmountValue }}" currencyCode="{{ .AmountCurrencyCode }}" exponent="{{ .AmountExponent }}"/>
          <orderContent />
          <paymentDetails>
            <APPLEPAY-SSL>
              <header>
                <ephemeralPublicKey>{{ .Payment.Token.PaymentData.Header.EphemeralPublicKey }}</ephemeralPublicKey>
                <publicKeyHash>{{ .Payment.Token.PaymentData.Header.PublicKeyHash }}</publicKeyHash>
                <transactionId>{{ .Payment.Token.PaymentData.Header.TransactionID }}</transactionId>
              </header>
              <signature>{{ .Payment.Token.PaymentData.Signature }}</signature>
              <version>{{ .Payment.Token.PaymentData.Version }}</version>
              <data>{{ .Payment.Token.PaymentData.Data }}</data>
            </APPLEPAY-SSL>
          </paymentDetails>
          <shopper>
            <shopperEmailAddress>{{ .Payment.ShippingContact.EmailAddress }}</shopperEmailAddress>
          </shopper>
        </order>
      </submit>
    </paymentService>
    `)
}

func wapProcess(merchantcode string, password string, payload json.RawMessage) string {

	var trans wapTransaction
	err := json.Unmarshal(payload, &trans)
	if err != nil {
		log.Printf("error parsing wap transaction %v", err)
		return "error parsing wap transaction"
	}

	// decode payment object (base 64 encoded JSON)
	paymentJSON, err := base64.StdEncoding.DecodeString(trans.Payment)
	if err != nil {
		log.Printf("error decoding payment object %v", err)
		return "error decoding payment object"
	}
	if len(paymentJSON) > 0 {
		err = json.Unmarshal(paymentJSON, &trans.PaymentObject)
		if err != nil {
			log.Printf("error parsing payment object %v", err)
			return "error parsing payment object"
		}
	}

	if wapTemplateErr != nil {
		log.Printf("error parsing wap template %v", wapTemplateErr)
		return "error parsing wap template"
	}

	trans.MerchantCode = merchantcode
	log.Printf("wap transaction:\n %v", trans)

	var wapRequest bytes.Buffer
	err = wapTemplate.Execute(&wapRequest, trans)
	if err != nil {
		log.Printf("error executing wap template %v", err)
		return "error executing wap template"
	}

	log.Printf("wap request:\n %s", string(wapRequest.Bytes()))

	req, err := http.NewRequest("POST", "https://secure-test.worldpay.com/jsp/merchant/xml/paymentService.jsp", &wapRequest)
	if err != nil {
		log.Printf("error preparing wap request %v", err)
		return "error preparing wap request"
	}
	req.Header.Set("Content-Type", "text/xml")
	req.SetBasicAuth(merchantcode, password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error transporting  wap request %v", err)
		return "error transporting  wap request"
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("error status when transporting wap request %s", res.Status)
		return "error status when transporting wap request"
	}

	// Defer closing of underlying connection so it can be re-used
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error returning wap response %v", err)
		return "error returning wap response"
	}

	log.Printf("wap response:\n %v", res)
	log.Printf("wap response body:\n %s", string(body))

	return "TODO"
}
