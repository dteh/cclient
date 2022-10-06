package cclient

import (
	"golang.org/x/net/proxy"

	http "github.com/dteh/fhttp"
	utls "gitlab.com/yawning/utls.git"
)

func NewClient(clientHello utls.ClientHelloID, customClientHello []byte, proxyUrl ...string) (http.Client, error) {
	var gen *utls.ClientHelloSpec
	var err error
	if clientHello == utls.HelloCustom && len(customClientHello) > 0 {
		f := utls.Fingerprinter{}
		gen, err = f.FingerprintClientHello(customClientHello)
		if err != nil {
			return http.Client{}, err
		}
	}

	if len(proxyUrl) > 0 && len(proxyUrl[0]) > 0 {
		dialer, err := newConnectDialer(proxyUrl[0])
		if err != nil {
			return http.Client{}, err
		}
		return http.Client{
			Transport: newRoundTripper(clientHello, gen, dialer),
		}, nil
	} else {
		return http.Client{
			Transport: newRoundTripper(clientHello, gen, proxy.Direct),
		}, nil
	}
}
