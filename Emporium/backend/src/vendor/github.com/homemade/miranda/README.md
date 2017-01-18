# Miranda [![GoDoc](https://godoc.org/github.com/homemade/miranda?status.svg)](https://godoc.org/github.com/homemade/miranda)

`miranda` is a Go package providing a server side dance partner for Apple Pay Merchant Validation.

It is named after Carmen Miranda, the dancer in the tutti-frutti hat.

## Overview

In order to accept Apple Pay payments on the web a merchant session object is required.

The `miranda` package provides a Dance function to perform the server side steps through it's MerchantValidationService interface.

A partner is required for this dance in order to perform the client side steps and provide the validation URL.

See https://developer.apple.com/reference/applepayjs/applepaysession#2166532
