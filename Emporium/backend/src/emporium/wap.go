package main

import "encoding/json"

var wapTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE paymentService PUBLIC "-//WorldPay/DTD WorldPay PaymentService v1//EN"
http://dtd.worldpay.com/paymentService_v1.dtd">
<paymentService version="1.4" merchantCode="{{ .MerchantCode }}">
<submit>
<order orderCode="{{ .OrderCode }}" shopperLanguageCode="{{ .ShopperLanguageCode }}"
<description>{{ .Description }}</description>
<amount value="{{ .Value }}" currencyCode="{{ .CurrencyCode }}" exponent="{{ .Exponent }}"/>
<
orderContent>
<![CDATA[]]>
</orderContent>
<paymentDetails>
<APPLEPAY
-
SSL>
<header>
<ephemeralPublicKey>{{ .EphemeralPublicKey }}</ephemeralPublicKey>
<publicKeyHash>{{ .PublicKeyHash }}</publicKeyHash>
<transactionId>{{ .TransactionId }}</transactionId>
</header>
<signature>{{ .Signature }}</signature>
<version>EC_v1</version>
<data>{{ .Data }}</data>
</APPLEPAY-SSL>
</paymentDetails>
<shopper>
<shopperEmailAddress>{{ .ShopperEmailAddress }}</shopperEmailAddress>
</shopper>
</order>
</submit>
</paymentService
`

func wapProcess(payment json.RawMessage) string {
	return "TODO"
}
