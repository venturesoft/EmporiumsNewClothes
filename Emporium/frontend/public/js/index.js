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
		requiredBillingContactFields: [
						"postalAddress",
						"name"
		],
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
		if (event.payment) {
			// create a random order code for testing purposes...
			var rnd = function() {
				return Math.floor((1 + Math.random()) * 0x10000)
					.toString(16)
					.substring(1);
			}
			var orderCode = rnd() + "-" + rnd() + "-" + rnd();

			// if we have some payment data add it
			var paymentDataEnc = "";
			if (event.payment.token && event.payment.token.paymentData) {
				// we send this as base64 encoded json
				paymentDataEnc = window.btoa(JSON.stringify(event.payment.token.paymentData));
			}

			// add the shopper's email address if we have it
			var shopperEmailAddress = "";
			if (event.payment.shippingContact && event.payment.shippingContact.emailAddress) {
				shopperEmailAddress = event.payment.shippingContact.emailAddress;
			}

			var transaction = {
				PaymentDataEnc: paymentDataEnc,
				OrderCode: orderCode,
				OrderDescription: "Donation",
				ShopperLanguageCode: "en",
				ShopperEmailAddress: shopperEmailAddress,
				AmountValue: "199",
				AmountCurrencyCode: "GBP",
				AmountExponent: "2"
			}
			// Send payment for processing...
			processPayment(transaction).then(function(result) {
				// ...return a status and redirect to a confirmation page
				session.completePayment(ApplePaySession.STATUS_SUCCESS);
				window.location.href = "/success.html?order=" + orderCode;
			});
		} else {
			alert("Unexpected payment event: " + JSON.stringify(event));
		}
	}

	// All our handlers are setup - start the Apple Pay payment
	session.begin();
}
