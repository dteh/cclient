package cclient

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"testing"
	"time"

	fhttp "github.com/dteh/fhttp"
	tls "github.com/refraction-networking/utls"
)

func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
	_, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(resp))
	w.Write([]byte("OK"))
}

func TestGetTransportRaceCondition(t *testing.T) {
	// RoundTripper can panic with multiple calls to getTransport
	server := http.Server{
		Addr:    ":443",
		Handler: http.HandlerFunc(httpRequestHandler),
	}
	defer server.Close()
	go server.ListenAndServeTLS("./keys/cert.pem", "./keys/private-key.pem")

	// This will trigger a panic if no lock is implemented in RoundTrip
	numReqs := 100
	wg := sync.WaitGroup{}
	wg.Add(numReqs)

	cl, _ := NewClient(tls.HelloChrome_Auto, &NewClientSettings{InsecureSkipVerify: true})
	f := func(cl *fhttp.Client) {
		resp, err := cl.Get("https://localhost")
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		wg.Done()
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		go f(&cl)
	}
	wg.Wait()

	log.Println("took:", time.Since(start))
}

func httpRequestHandlerDontClose(w http.ResponseWriter, req *http.Request) {
	_, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(resp))
	w.WriteHeader(150)
	// <-make(chan bool)
	w.Write([]byte("OK"))
}
func TestBodyTimeoutLeak(t *testing.T) {
	server := http.Server{
		Addr:    ":443",
		Handler: http.HandlerFunc(httpRequestHandlerDontClose),
	}
	defer server.Close()
	go func() error {
		ln, err := net.Listen("tcp", ":443")
		if err != nil {
			return err
		}
		defer ln.Close()
		return server.ServeTLS(ln, "./keys/cert.pem", "./keys/private-key.pem")
	}()

	cl, _ := NewClient(tls.HelloChrome_Auto, &NewClientSettings{InsecureSkipVerify: true}, "http://127.0.0.1:8888")
	cl.Timeout = 2 * time.Second
	r, err := cl.Get("https://localhost")
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(b) + "S")
}
