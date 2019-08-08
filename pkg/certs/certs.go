package certs

// generate-tls-cert generates root, leaf, and client TLS certificates.

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/plunder-app/plunder/pkg/utils"
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

	// Find all IP addresses on a server
	serverAddresses, err := utils.FindAllIPAddresses()
	if err != nil {
		return err
	}

	// Find the hostname of the server
	serverName, err := os.Hostname()
	if err != nil {
		return err
	}

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
		IPAddresses:  serverAddresses,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{serverName},
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

	keyData = certPrivKeyPEM.Bytes()

	return nil
}

// WriteKeyToFile - will write a generated Key to a file path
func WriteKeyToFile(path string) error {

	err := ioutil.WriteFile(path, keyData, 0600)
	if err != nil {
		return err
	}
	return nil
}

// WritePemToFile - will write a generated pem to a file path
func WritePemToFile(path string) error {

	err := ioutil.WriteFile(path, pemData, 0600)
	if err != nil {
		return err
	}

	return nil
}

// GetKey - will return the []byte of the key
func GetKey() []byte {
	return keyData
}

// GetPem - will return the []byte of the key
func GetPem() []byte {
	return pemData
}
