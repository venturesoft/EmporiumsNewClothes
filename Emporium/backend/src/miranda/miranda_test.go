package miranda

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testRequest = `
{
	"merchantIdentifier": "merchant.com.example",
	"displayName": "Development",
	"domainName": "example.com"
}
`

type TestMerchantValidationConfig struct {
}

func (c TestMerchantValidationConfig) LoadCert() (tls.Certificate, error) {
	testPair := []byte(`
Bag Attributes
    localKeyID: 66 B4 11 E0 DD 40 85 93 01 11 30 FB 8A 75 D3 C6 27 F9 03 A2
subject=/C=GB/ST=Some-State/O=Internet Widgits Pty Ltd
issuer=/C=GB/ST=Some-State/O=Internet Widgits Pty Ltd
-----BEGIN CERTIFICATE-----
MIIFtTCCA52gAwIBAgIJALL6Zs/4J/F3MA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkdCMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTcwMTExMDAwNTExWhcNMzcwMTA2MDAwNTExWjBF
MQswCQYDVQQGEwJHQjETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAwluyg2UoFdDSsIroPNhWZk+2IWVBmRsSPOaBftdbwuEiWg3ze1VN2fim
KQWQHz1ex4htvvXiTad56Vq/UftK/Khvt8OfF5TPZNgykDaSJwt2jrfujRh+oxrq
Hwu5N3h4/ERxjqLaDf+D+m7yyHV3BolkqpLw6JTlm2Kwej7fOD4dzzMu4h+lO12f
FZpEof50PGb998WUhPeztv21ly9nWfwygd8QIkGnq6PkGFcGfqZYR+JKqPbWQSp3
p436+TIfnul5O96JS8P5zdeduMh1kUXXdf83FijrVKJT7XUqOoKm9dhfX0no9v10
ewBSqoA6cxN1CsiyK5UCSIvIiGyBOmiMCwP7jHkGJ7lCNG0oLZo40Ux1O1oK50Li
lcT0Agj2jglHLGFuK4LvT3QB0aVJhCOVgnph8dTyHQ8eahFuBxKmK+aNXM1KXcZr
PFdid5f3o7M0eZz/9ce3arIR/oP7PYAcmkGxTZ9QckYxUsIAd3DChH9FZfASiuAM
RPH7IxgK1+xE4R6eIqjyL/+DzSStA41jB7WHECXsORxoDM6Eg3NXwDKqHFqn/CCF
R356hGkOzyGFd8QRhrPqB3e9z8G/TSWij78siMIJJNBxMPzIofCVyjQfk2mIwgIl
G81PhAaNo9+0ExtFpxDHHlWDB0801TzbJXMiCVT2Khr5wYxDMPsCAwEAAaOBpzCB
pDAdBgNVHQ4EFgQUruXmX7zP2YZo2lyJDEXycKJ3ZkwwdQYDVR0jBG4wbIAUruXm
X7zP2YZo2lyJDEXycKJ3ZkyhSaRHMEUxCzAJBgNVBAYTAkdCMRMwEQYDVQQIEwpT
b21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGSCCQCy
+mbP+CfxdzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBBQUAA4ICAQC0IujZS7l/
Gv1JLl12GrHy+NDGl0ja2jAtl0niboP2kzimw3ui36dx8fpcm/4HMvo1z0cNxyip
oOnYic27RheqRLAH97C8tw5mYPcAP/OvJv+yWmsb8jgOH33NK25ZIyJYBZckXbxo
V4GrCm00QIjhoGrhmC1Nz+2+wqn/TBn5ww6JuNgbwrXQe57WdAjkN2LYXbo6quQv
JiGGHlmiPvZTWgSHtiyuoXC6eHiflI1WXtJbN8ziJpSjKAbetx9gLRaKZ4F4U9D/
tUPgLweS0t+cR+jYw2rz9oqnyWNXlq+rwmb7VMINTpX9BI5dpfOeV3UiLygnz/W+
jsqBS1V0ewnrr+V4KNwYZZ6Yp3drSICo1rgfp0IT/aW5G9dEDQIkJm1DKIlopHzg
/LUjMX0IezRZH/8n73l6PYpd+EfhLz12U7ru26Tjtr3aIoZMrGUAyX8AfT4RaEtN
4FLErbflG9+F37WrfCDByN9TZvQd/FK5yttJSAVh3QN7r2Y0GtQywlaDQvbzJnKS
pOcY63p6Fs+EjBns3VSmJ13EmFZIH1JjuSrxas+nrNwrXShHZsG8FKetKhkiGrOK
B4xdQEBj5VYUb9u7FdvJSkODt0MIOsvNzshrDcOCi3s4hYPSNuSksMN15J0UMC1b
FtURHTw+dsBT1xHKGm9IoUzchOSifZSzSA==
-----END CERTIFICATE-----
Bag Attributes
    localKeyID: 66 B4 11 E0 DD 40 85 93 01 11 30 FB 8A 75 D3 C6 27 F9 03 A2
Key Attributes: <No Attributes>
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)
	return tls.X509KeyPair(testPair, testPair)
}
func (c TestMerchantValidationConfig) LoadRequestBody() (json.RawMessage, error) {
	return []byte(testRequest), nil
}

func TestMiranda(t *testing.T) {

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()

	carmen, err := CreateMerchantValidationService(30*time.Second, TestMerchantValidationConfig{})
	if err != nil {
		t.Fatal(err)
	}

	session, err := carmen.Dance(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if string(session) != "OK" {
		t.Fatal(fmt.Errorf("invalid session %s", string(session)))
	}

}
