package encode

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	// RSAPrivateKeyBlockType is a possible value for pem.Block.Type.
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
	// CertificateBlockType is a possible value for pem.Block.Type.
	CertificateBlockType = "CERTIFICATE"
)

// CertPEM returns PEM-encoded certificate data
func CertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// PrivateKeyPEM returns PEM-encoded private key data
func PrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}
