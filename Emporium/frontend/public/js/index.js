/*
Copyright (C) 2016 Apple Inc. All Rights Reserved.
See LICENSE.txt for this sample’s licensing information

Abstract:
The main client-side JS. Handles displaying the Apple Pay button and requesting a payment.
*/

/**
* This method is called when the page is loaded.
* We use it to show the Apple Pay button as appropriate.
* Here we're using the ApplePaySession.canMakePayments() method,
* which performs a basic hardware check.
*
* If we wanted more fine-grained control, we could use
* ApplePaySession.canMakePaymentsWithActiveCards() instead.
*/
document.addEventListener('DOMContentLoaded', () => {
	if (window.ApplePaySession) {
		if (ApplePaySession.canMakePayments) {
			showApplePayButton();
		}
	}
});

function showApplePayButton() {
	HTMLCollection.prototype[Symbol.iterator] = Array.prototype[Symbol.iterator];
	const buttons = document.getElementsByClassName("apple-pay-button");
	for (let button of buttons) {
		button.className += " visible";
	}
}


/**
* Apple Pay Logic
* Our entry point for Apple Pay interactions.
* Triggered when the Apple Pay button is pressed
*/
function applePayButtonClicked() {
	const paymentRequest = {
		countryCode: 'GB',
		currencyCode: 'GBP',
		total: {
			label: 'Apple Pay Example',
			amount: '1.99',
		},
		supportedNetworks:[ 'amex', 'masterCard', 'visa'],
		merchantCapabilities: [ 'supports3DS' ],
		requiredShippingContactFields: [ 'email' ],
	};

	const session = new ApplePaySession(1, paymentRequest);

	/**
	* Merchant Validation
	* We call our merchant session endpoint, passing the URL to use
	*/
	session.onvalidatemerchant = (event) => {
		console.log("Validate merchant");
		const validationURL = event.validationURL;
		getApplePaySession(event.validationURL).then(function(response) {
  			console.log(response);
  			session.completeMerchantValidation(response);
		});
	};

	/**
	* Shipping Method Selection
	* If the user changes their chosen shipping method we need to recalculate
	* the total price. We can use the shipping method identifier to determine
	* which method was selected.
	*/
	session.onshippingmethodselected = (event) => {
		session.completeShippingMethodSelection(ApplePaySession.STATUS_SUCCESS, total, lineItems);
	};

	/**
	* Payment Authorization
	* Here you receive the encrypted payment data. You would then send it
	* on to your payment provider for processing, and return an appropriate
	* status in session.completePayment()
	*/
	session.onpaymentauthorized = (event) => {
		// Send payment for processing...
		const payment = event.payment;

		// ...return a status and display the result
		session.completePayment(ApplePaySession.STATUS_SUCCESS);
		alert(JSON.stringify(payment));

	}

	// All our handlers are setup - start the Apple Pay payment
	session.begin();
}
