package kcpproxy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Server holds the configuration for the proxy server
type Server struct {
	// The hostname:port for the proxy to listen on
	ListenAddress string

	// CA used to validate client certs connecting to the proxy
	ClientCACert string

	// The proxy's server cert
	ServerCertFile string

	// The proxy's private key file
	ServerKeyFile string

	// A yaml file containing a list of Paths
	MappingFile string
}

// PathMapping describes how to route traffic from a path to a backend server.
// The yaml file is a list of these objects. Path is registered in the
// DefaultServeMux with a handler that delegates to the backend.
type PathMapping struct {
	Path            string `yaml:"path"`
	Backend         string `yaml:"backend"`
	BackendServerCA string `yaml:"backend_server_ca"`
	ProxyClientCert string `yaml:"proxy_client_cert"`
	ProxyClientKey  string `yaml:"proxy_client_key"`
	UserHeader      string `yaml:"user_header,omitempty"`
	GroupHeader     string `yaml:"group_header,omitempty"`
}

func (s *Server) Serve() error {
	mappingData, err := ioutil.ReadFile(s.MappingFile)
	if err != nil {
		return err
	}

	mapping := []PathMapping{}
	if err = yaml.Unmarshal(mappingData, &mapping); err != nil {
		return err
	}

	for _, pathCfg := range mapping {
		proxy, err := NewReverseProxy(pathCfg.Backend, pathCfg.ProxyClientCert, pathCfg.ProxyClientKey, pathCfg.BackendServerCA)
		if err != nil {
			return err
		}
		userHeader := "X-Remote-User"
		groupHeader := "X-Remote-Group"
		if pathCfg.UserHeader != "" {
			userHeader = pathCfg.UserHeader
		}
		if pathCfg.GroupHeader != "" {
			groupHeader = pathCfg.GroupHeader
		}
		http.Handle(pathCfg.Path, http.HandlerFunc(ProxyHandler(proxy, userHeader, groupHeader)))
	}

	clientCACert, err := ioutil.ReadFile(s.ClientCACert)
	if err != nil {
		return err
	}

	clientCACertPool := x509.NewCertPool()
	clientCACertPool.AppendCertsFromPEM(clientCACert)

	server := &http.Server{
		Addr:    s.ListenAddress,
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.VerifyClientCertIfGiven,
			ClientCAs:  clientCACertPool,
		},
	}

	return server.ListenAndServeTLS(s.ServerCertFile, s.ServerKeyFile)
}
