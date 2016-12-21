package main

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"
)

type wapTransaction struct {
	Payment struct {
		Token struct {
			PaymentData struct {
				Data      string `json:"data"`
				Signature string `json:"signature"`
				Header    struct {
					PublicKeyHash      string `json:"publicKeyHash"`
					EphemeralPublicKey string `json:"ephemeralPublicKey"`
					TransactionId      string `json:"transactionId"`
				} `json:"header"`
				Version string `json:"version"`
			} `json:"paymentData"`
			TransactionId string `json:"transactionId"`
			PaymentMethod struct {
				Network     string `json:"network"`
				Type        string `json:"type"`
				DisplayName string `json:"displayName"`
			} `json:"paymentMethod"`
		} `json:"token"`
		ShippingContact struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"shippingContact"`
	}
	OrderCode           string
	OrderDescription    string
	ShopperLanguageCode string
	AmountValue         string
	AmountCurrencyCode  string
	AmountExponent      string
	merchantCode        string
}

var wapTemplate *template.Template
var wapTemplateErr error

func init() {
	wapTemplate, wapTemplateErr = template.New("wapTemplate").Parse(`<?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE paymentService PUBLIC "-//WorldPay/DTD WorldPay PaymentService v1//EN"
    http://dtd.worldpay.com/paymentService_v1.dtd">
    <paymentService version="1.4" merchantCode="{{ .MerchantCode }}">
    <submit>
    <order orderCode="{{ .OrderCode }}" shopperLanguageCode="{{ .ShopperLanguageCode }}"
    <description>{{ .OrderDescription }}</description>
    <amount value="{{ .AmountValue }}" currencyCode="{{ .AmountCurrencyCode }}" exponent="{{ .AmountExponent }}"/>
    <orderContent>
    <![CDATA[]]>
    </orderContent>
    <paymentDetails>
    <APPLEPAY-SSL>
    <header>
    <ephemeralPublicKey>{{ .Token.PaymentData.Header.EphemeralPublicKey }}</ephemeralPublicKey>
    <publicKeyHash>{{ .Token.PaymentData.Header.PublicKeyHash }}</publicKeyHash>
    <transactionId>{{ .Token.PaymentData.Header.TransactionId }}</transactionId>
    </header>
    <signature>{{ .Token.PaymentData.Signature }}</signature>
    <version>{{ .Token.PaymentData.Version }}</version>
    <data>{{ .Token.PaymentData.Data }}</data>
    </APPLEPAY-SSL>
    </paymentDetails>
    <shopper>
    <shopperEmailAddress>{{ .ShippingContact.EmailAddress }}</shopperEmailAddress>
    </shopper>
    </order>
    </submit>
    </paymentService
    `)
}

func wapProcess(merchantcode string, payload json.RawMessage) string {
	log.Printf("payload:\n %s", string(payload))

	var trans wapTransaction
	err := json.Unmarshal(payload, &trans)
	if err != nil {
		log.Printf("error parsing wap transaction %v", err)
		return "error parsing wap transaction"
	}

	if wapTemplateErr != nil {
		log.Printf("error parsing wap template %v", wapTemplateErr)
		return "error parsing wap template"
	}

	trans.merchantCode = merchantcode
	log.Printf("wap transaction:\n %v", trans)

	var wapRequest bytes.Buffer
	err = wapTemplate.Execute(&wapRequest, trans)
	if err != nil {
		log.Printf("error executing wap template %v", err)
		return "error executing wap template"
	}

	log.Printf("wap request:\n %s", string(wapRequest.Bytes()))

	return "TODO"
}
