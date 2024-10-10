package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

func main() {

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}

	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		RootCAs:      caCertPool,
	}

	server := &http.Server{
		Addr:      ":3443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		certs := r.TLS.PeerCertificates
		fmt.Println(certs[0].Subject.CommonName)
		w.Write([]byte("Hello, secure world!"))
	})

	// http.ListenAndServe(":3000", nil)
	// /etc/hosts
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
