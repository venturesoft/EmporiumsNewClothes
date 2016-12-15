# Emporium Web's New Clothes
Thsi is a fork of Apple's [Emporium Web](https://developer.apple.com/library/content/samplecode/EmporiumWeb/Introduction/Intro.html) simple one page site designed to show you how to request and handle Apple Pay payments on the web. 

Apple Pay on the web consists of both a client-side and server-side component. You request payment on the client, and validate yourself as a merchant on your web server. This ~~node.js~~ Go example shows you how to carry out both the payment request and the merchant validation.

## Requirements
This example is a self-contained project that uses ~~node.js and Express~~ Docker to run a small web server. 

## Getting Started

#### Generate your Apple Pay Certificates
Apple Pay requires a merchant identifier and two certificates - a *session* certificate and a *rewrap* certificate. The merchant identifier uniquely identifies you as an Apple Pay merchant. The *rewrap* certificate is used to encrypt your Apple Pay payments, and the *session* certificate is used to authenticate your website.

Create your Apple Pay merchant identifier at https://developer.apple.com, and register your web domain against it. Convert your session certificate and key to `PEM` format, and place it in this example's `Emporium/certificates` directory.

> NOTE
> Session certificates must be imported to the `login` keychain so that the private key is accessible (they can then be exported in P12 format). Final conversion to PEM format can be achieved using OpenSSL

    openssl pkcs12 -in cert.p12 -out cert.pem -nodes -clcerts

#### Set up SSL
Apple Pay requires your site to be hosted over HTTPS. Generate your SSL certificate, and place the certificate and key in this example's `Emporium/certificates` directory.

#### Run the example    
```
git clone https://github.com/venturesoft/emporiumsnewclothes && cd emporiumsnewclothes
printf 'PRIVATE_DIR=%s\n' /home/user/private > .env
docker-compose build && docker-compose up
```

## Resources
A number of resources are available to help you with Apple Pay. 

  * Apple Pay Developer Site - https://developer.apple.com/apple-pay/
  * Apple Pay on the web WWDC Session Video - https://developer.apple.com/videos/play/wwdc2016/703/
  * Apple Pay Domain Verification - https://developer.apple.com/support/apple-pay-domain-verification/

## ES6
The client-side code for this example is written in ES6 (the latest version of Safari supports ES6 natively).

Original source Copyright (C) 2016 Apple Inc. All rights reserved.
