package certs

// generate-tls-cert generates root, leaf, and client TLS certificates.

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

// Internal variables to hold the outputs
var keyData, pemData []byte

// GenerateKeyPair - (TODO)
func GenerateKeyPair(hosts []string, start time.Time, length time.Duration) error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Plunder"},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"Yorkshire"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Plunder"},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"Yorkshire"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		//DNSNames:     []string{"deploy01"}, // TODO
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	pemData = certPEM.Bytes()

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	_, err = tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		fmt.Printf("Da fuck [%s]\n", err.Error())
	}

	_, err = tls.LoadX509KeyPair("plunder.pem", "plunder.key")
	if err != nil {
		fmt.Printf("Da fuck [%s]\n", err.Error())
	}

	keyData = certPrivKeyPEM.Bytes()

	// serverTLSConf = &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// }

	// certpool := x509.NewCertPool()
	// certpool.AppendCertsFromPEM(caPEM.Bytes())
	// clientTLSConf = &tls.Config{
	// 	RootCAs: certpool,
	// }

	return nil
}

// WriteKeyToFile - will write a generated Key to a file path
func WriteKeyToFile(path string) error {

	err := ioutil.WriteFile(path, keyData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// WritePemToFile - will write a generated pem to a file path
func WritePemToFile(path string) error {

	err := ioutil.WriteFile(path, pemData, 0644)
	if err != nil {
		return err
	}

	return nil
}
