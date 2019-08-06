package certs

// generate-tls-cert generates root, leaf, and client TLS certificates.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// Internal variables to hold the outputs
var keyData, pemData []byte

// GenerateKeyPair - (TODO)
func GenerateKeyPair(hosts []string, start time.Time, length time.Duration) error {
	// Sanity check inputs

	// Hosts will need checkign at somepoint (TODO)

	//	if len(hosts) == 0 {
	//		return fmt.Errorf("No hosts have been submitted")
	//	}

	notAfter := start.Add(length)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	keyData, err = x509.MarshalECPrivateKey(rootKey)
	if err != nil {
		return fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
	}

	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Plunder"},
			CommonName:   "Plunder CA",
		},
		NotBefore:             start,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	pemData, err = x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		return err
	}

	return nil
}

// WriteKeyToFile - will write a generated Key to a file path
func WriteKeyToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyData}); err != nil {
		return err
	}
	return nil
}

// WritePemToFile - will write a generated pem to a file path
func WritePemToFile(path string) error {
	certOut, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %s", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: pemData}); err != nil {
		return fmt.Errorf("failed to write data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing cert.pem: %s", err)
	}
	return nil
}
