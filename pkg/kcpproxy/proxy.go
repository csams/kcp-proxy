package kcpproxy

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"crypto/x509"
)

// NewReverseProxy returns a new reverse proxy where backend is the backend URL to
// connect to, clientCert is the proxy's client cert to use to connect to it,
// clientKeyFile is the proxy's client private key file, and caFile is the CA
// the proxy uses to verify the backend server's cert.
func NewReverseProxy(backend, clientCert, clientKeyFile, caFile string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(backend)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(clientCert, clientKeyFile)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	return proxy, nil
}

// ProxyHandler extracts the CN as a user name and Organizations as groups from
// the client cert if one is provided and adds them as HTTP headers to the request
// that gets forwarded to the backend
func ProxyHandler(proxy *httputil.ReverseProxy, UserHeader, GroupHeader string) func(wr http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) >= 1 {
			clientCert := r.TLS.PeerCertificates[0]
			appendClientCertAuthHeaders(r.Header, clientCert, UserHeader, GroupHeader)
		}
		proxy.ServeHTTP(w, r)
	}
}

func appendClientCertAuthHeaders(header http.Header, clientCert *x509.Certificate, UserHeader, GroupHeader string) {
	userName := clientCert.Subject.CommonName
	header.Set(UserHeader, userName)

	groups := clientCert.Subject.Organization
	for _, group := range groups {
		header.Add(GroupHeader, group)
	}
}
