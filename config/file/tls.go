package file

import "crypto/tls"

type TLS struct {
	Certificates []Certificate `yaml:"certificates"`
}

type Certificate struct {
	CertFile string `yaml:"certfile"`
	KeyFile  string `yaml:"keyfile"`
}

func (h TLS) GetConfig() (*tls.Config, error) {
	if len(h.Certificates) == 0 {
		return nil, nil
	}

	certs := make([]tls.Certificate, 0)
	for _, c := range h.Certificates {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return &tls.Config{
		Certificates: certs,
		NextProtos:   []string{"h2"},
	}, nil
}
