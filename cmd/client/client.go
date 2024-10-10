package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("robson.crt", "robson.key")
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
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	res, err := client.Get("https://localhost:3443/")
	if err != nil {
		log.Fatal(err)
	}
	res.Write(os.Stdout)
	defer res.Body.Close()
}
