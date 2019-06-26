package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math/big"
	"time"

	"github.com/d-kuro/sealed-secrets-key/pkg/encode"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Generate(name string, namespace string, keySize int, keyTTL time.Duration) ([]byte, error) {
	r := rand.Reader
	key, err := rsa.GenerateKey(r, keySize)
	if err != nil {
		return nil, err
	}
	cert, err := signKey(r, key, keyTTL)
	if err != nil {
		return nil, err
	}
	certs := []*x509.Certificate{cert}
	s, err := createSecret(key, certs, namespace, name)
	if err != nil {
		return nil, err
	}
	return s.AsYAML()
}

func signKey(r io.Reader, key *rsa.PrivateKey, keyTTL time.Duration) (*x509.Certificate, error) {
	notBefore := time.Now()

	serialNo, err := rand.Int(r, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	cert := x509.Certificate{
		SerialNumber: serialNo,
		KeyUsage:     x509.KeyUsageEncipherOnly,
		NotBefore:    notBefore.UTC(),
		NotAfter:     notBefore.Add(keyTTL).UTC(),
		Subject: pkix.Name{
			CommonName: "",
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	data, err := x509.CreateCertificate(r, &cert, &cert, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(data)
}

func createSecret(key *rsa.PrivateKey, certs []*x509.Certificate, namespace, keyName string) (*Secret, error) {
	certBytes := make([]byte, 0, len(certs))
	for _, cert := range certs {
		certBytes = append(certBytes, encode.CertPEM(cert)...)
	}

	s := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              keyName,
			Namespace:         namespace,
			CreationTimestamp: metav1.Time{},
		},
		Data: map[string][]byte{
			v1.TLSPrivateKeyKey: encode.PrivateKeyPEM(key),
			v1.TLSCertKey:       certBytes,
		},
		Type: v1.SecretTypeTLS,
	}
	s.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Secret"))

	return CopyKubernetesSecret(s)
}
