package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	microerror "github.com/giantswarm/microkit/error"
)

type CertFiles struct {
	RootCAs []string // Root certificate authority file paths.
	Cert    string   // X.509 certificate file path.
	Key     string   // X.509 key file path.
}

// LoadTLSConfig creates TLS configuration for given crtificate files. It
// assumes X.509 keypair and sets minimum 1.2 minimum TLS version. All fields
// of CertFiles are optional. If the field is missing, the corresponding
// certificate will not be loaded.
func LoadTLSConfig(files CertFiles) (*tls.Config, error) {
	var (
		loadCert    = files.Cert != "" && files.Key != ""
		loadRootCAs = len(files.RootCAs) > 0
	)

	if !loadCert && !loadRootCAs {
		return nil, nil
	}

	var certificate tls.Certificate
	if loadCert {
		cert, err := ioutil.ReadFile(files.Cert)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		key, err := ioutil.ReadFile(files.Key)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		certificate, err = tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var rootCAs *x509.CertPool
	if loadRootCAs {
		rootCAs = x509.NewCertPool()

		for _, caFile := range files.RootCAs {
			pemByte, err := ioutil.ReadFile(caFile)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			for {
				var block *pem.Block
				block, pemByte = pem.Decode(pemByte)
				if block == nil {
					break
				}
				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				rootCAs.AddCert(cert)
			}
		}
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      rootCAs,
		MinVersion:   tls.VersionTLS12,
	}
	return &tlsConfig, nil
}
