package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

type CertInfo struct {
	Organization string
	Country      string
	Locality     string
	Domain       string
}

func CreateCA() error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization:       []string{"CA Corp"},
			Country:            []string{"BR"},
			Locality:           []string{"NATAL"},
			CommonName:         "RobsonDev CA",
			OrganizationalUnit: []string{"TI"},
			Province:           []string{"RN"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	//gerar a chave privada
	caPK, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	//gerar a chave pública ou o certificado
	caData, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPK.PublicKey, caPK)
	if err != nil {
		return err
	}

	// escrever a chave pública em arquivos no formato PEM
	caPEM, err := os.Create("ca.crt")
	if err != nil {
		return err
	}
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caData,
	})

	//escrever a chave privada no formato PEM
	caPrivateKeyPEM, err := os.Create("ca.key")
	if err != nil {
		return err
	}
	pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPK),
	})

	return nil
}

func CreateAndSignCertificate(caFileName string, certName string, info CertInfo) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{info.Organization},
			Country:      []string{info.Country},
			Locality:     []string{info.Locality},
			CommonName:   info.Domain,
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{info.Domain, "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPK, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	caPK, ca, err := loadCertificates(caFileName)
	if err != nil {
		return err
	}
	certData, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPK.PublicKey, caPK)
	if err != nil {
		return err
	}

	certPEM, err := os.Create(certName + ".crt")
	if err != nil {
		return err
	}
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certData,
	})

	certPrivKeyPEM, err := os.Create(certName + ".key")
	if err != nil {
		return err
	}
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPK),
	})
	return nil
}

func loadCertificates(fileName string) (*rsa.PrivateKey, *x509.Certificate, error) {
	// Carregar chave privada da CA
	caPrivFile, err := os.ReadFile(fileName + ".key")
	if err != nil {
		return nil, nil, err
	}
	caPrivBlock, _ := pem.Decode(caPrivFile)
	caPriv, err := x509.ParsePKCS1PrivateKey(caPrivBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// Carregar certificado da CA
	caCertFile, err := os.ReadFile(fileName + ".crt")
	if err != nil {
		return nil, nil, err
	}
	caCertBlock, _ := pem.Decode(caCertFile)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return caPriv, caCert, nil
}

func ConvertToPKCS12(certName, caCertName, password string) error {
	_, caCert, err := loadCertificates(caCertName)
	if err != nil {
		return err
	}
	pk, cert, err := loadCertificates(certName)
	if err != nil {
		return err
	}
	data, err := pkcs12.Modern2023.Encode(pk, cert, []*x509.Certificate{caCert}, password)
	if err != nil {
		return err
	}
	err = os.WriteFile(certName+".p12", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
