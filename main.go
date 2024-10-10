package main

import "go-mtls/certificate"

func main() {
	// err := certificate.CreateCA()
	// if err != nil {
	// 	panic(err)
	// }

	// err = certificate.CreateAndSignCertificate("ca", "server", certificate.CertInfo{
	// 	Organization: "RobsonDev Corp",
	// 	Country:      "BR",
	// 	Locality:     "NATAL",
	// 	Domain:       "server.example.com",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// err = certificate.CreateAndSignCertificate("ca", "client", certificate.CertInfo{
	// 	Organization: "RobsonDev Corp",
	// 	Country:      "BR",
	// 	Locality:     "NATAL",
	// 	Domain:       "client.example.com",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// err := certificate.CreateAndSignCertificate("ca", "robson", certificate.CertInfo{
	// 	Organization: "RobsonDev Corp",
	// 	Country:      "BR",
	// 	Locality:     "NATAL",
	// 	Domain:       "robson.dev.br",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	err := certificate.ConvertToPKCS12("robson", "ca", "123456")
	if err != nil {
		panic(err)
	}
}
