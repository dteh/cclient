package cclient

import (
	"golang.org/x/net/proxy"

	http "github.com/dteh/fhttp"
	utls "github.com/refraction-networking/utls"
)

type NewClientSettings struct {
	CustomClientHello  []byte
	InsecureSkipVerify bool

	customClientHello *utls.ClientHelloSpec
}

func NewClient(clientHello utls.ClientHelloID, settings NewClientSettings, proxyUrl ...string) (http.Client, error) {
	var gen *utls.ClientHelloSpec
	var err error
	if clientHello == utls.HelloCustom && len(settings.CustomClientHello) > 0 {
		f := utls.Fingerprinter{}
		gen, err = f.FingerprintClientHello(settings.CustomClientHello)
		if err != nil {
			return http.Client{}, err
		}
		settings.customClientHello = gen
	}

	if len(proxyUrl) > 0 && len(proxyUrl[0]) > 0 {
		dialer, err := newConnectDialer(proxyUrl[0])
		if err != nil {
			return http.Client{}, err
		}
		return http.Client{
			Transport: newRoundTripper(clientHello, settings, dialer),
		}, nil
	} else {
		return http.Client{
			Transport: newRoundTripper(clientHello, settings, proxy.Direct),
		}, nil
	}
}
