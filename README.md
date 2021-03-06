# Emporium's New Clothes
This is a fork of Apple's [Emporium Web](https://developer.apple.com/library/content/samplecode/EmporiumWeb/Introduction/Intro.html) simple one page site designed to show you how to request and handle Apple Pay payments on the web (with the backend re-written to use Go).

Apple Pay on the web consists of both a client-side and server-side component. You request payment on the client, and validate yourself as a merchant on your web server. This ~~node.js~~ Go example shows you how to carry out both the payment request and the merchant validation.

## Requirements
This example is a self-contained project that uses ~~node.js and Express~~ Docker to run a small web server.

## Getting Started

#### Generate your Apple Pay Certificates
Apple Pay requires a merchant identifier and two certificates - a *session* certificate and a *rewrap* certificate. The merchant identifier uniquely identifies you as an Apple Pay merchant. The *rewrap* certificate is used to encrypt your Apple Pay payments, and the *session* certificate is used to authenticate your website.

Create your Apple Pay merchant identifier at https://developer.apple.com, and register your web domain against it. Convert your session certificate and key to `PEM` format, name it `merchant.pem`, then place it in a directory named `applepay` within a local directory for this project to read e.g. `/Users/darcy/private/EmporiumsNewClothes/applepay`.

> NOTE:
> When using a Mac, Session certificates must be imported to the `login` keychain so that the private key is accessible (they can then be exported in P12 format). Final conversion to PEM format can be achieved using OpenSSL e.g.

    openssl pkcs12 -in merchant.p12 -out merchant.pem -nodes -clcerts

#### Add a merchant configuration file
Add a file named `merchant.json` to the `applepay` local directory in the following format:
```
{
	"merchantIdentifier": "<Apple Merchant ID>",
	"displayName": "<Apple Merchant Display Name>",
	"domainName": "<Name Of Your Domain>"
}
```

#### Set up SSL
Apple Pay requires your site to be hosted over HTTPS. Generate your SSL certificate, and if required obtain intermediary certificates. Name the certificate bundle `bundle.crt` and name the key `private.key` and place them both in a directory named `ssl` within the same local directory used for the `applepay` directory e.g. `/Users/darcy/private/EmporiumsNewClothes/ssl`

#### Download the Apple domain association file
Apple uses a file to verify ownership of your domain e.g. `apple-developer-merchantid-domain-association`. Download this file and place it in a directory named `verification`, again in the same local directory e.g. `/Users/darcy/private/EmporiumsNewClothes/verification`

#### Run the example    
```
git clone https://github.com/venturesoft/EmporiumsNewClothes.git && cd EmporiumsNewClothes
printf 'PRIVATE_DIR=%s\n' /Users/darcy/private/EmporiumsNewClothes > .env && printf 'HTTP_PORT=%s\n' 80 >> .env && printf 'HTTPS_PORT=%s\n' 443 >> .env
docker-compose build && docker-compose up
```
> NOTE
> Replace `/Users/darcy/private/EmporiumsNewClothes` with the path to your local directory and if required specify different ports for HTTP/HTTPS

## Testing Apple Pay transactions
See the following guide for testing Apple Pay transactions - https://developer.apple.com/support/apple-pay-sandbox/

### Payment Processing
The default setup for the project will not process the Apple Pay transactions - the [payment](https://developer.apple.com/reference/applepayjs/payment) object will simply be logged in the backend.

To extend this default behaviour and submit the payment for processing by a Payment Gateway follow the appropriate providers instructions below:

#### Worldpay Corporate Gateway - XML Direct Integration

Add a file named `wap.json` to the `applepay` local directory in the following format:
```
{
	"merchantCode": "<Worldpay Merchant Code>",
	"password": "<Password for Worldpay XML Direct Integration>"
}
```

## Resources
A number of resources are available to help you with Apple Pay.

  * Apple Pay Developer Site - https://developer.apple.com/apple-pay/
  * Apple Pay on the web WWDC Session Video - https://developer.apple.com/videos/play/wwdc2016/703/
  * Apple Pay Domain Verification - https://developer.apple.com/support/apple-pay-domain-verification/

## ES6
The client-side code for this example is written in ES6 (the latest version of Safari supports ES6 natively).

Original source Copyright (C) 2016 Apple Inc. All rights reserved.
