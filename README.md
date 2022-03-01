# kcp-proxy
HTTPS proxy that maps paths to multiple backends, performs TLS re-encryption to
each, and supports client certificate based authentication that converts common
names and orgs to HTTP headers that are passed on to KCP.

# Example
The proxy itself needs a serving cert/key pair, and it needs a client cert and
key to identity itself to and communicate with each backend. It needs a CA for
its serving cert that clients connecting to it will need to trust, a CA that will
be used to verify all client certs used to connect to it, and the CA used to
trust the serving cert of each backend server.

KCP needs to be told which CA to use to trust the proxy's client cert. Start it
something like this:
```
#!/usr/bin/env bash

./bin/kcp start \
    ...
    --tls-cert-file=/home/<username>/projects/work/kcp/code/kcp/certs/server-cert.pem \
    --tls-private-key-file=/home/<username>/projects/work/kcp/code/kcp/certs/server-key.pem \
    --requestheader-client-ca-file=/home/<username>/projects/work/kcp/code/kcp-proxy/certs/ca-cert.pem \
    --requestheader-username-headers=X-Remote-User \
    --requestheader-group-headers=X-Remote-Group \
    ...
```

## Generate the proxy certs
Update the `.cnf` files in [hack](hack) so the certs include the IP address to
which KCP binds in their SANs. Then run this:
```bash
./hack/gen-certs.sh
```

## An example path mapping
```yaml
- path: /
  backend: https://localhost:6443
  backend_server_ca: /home/<username>/projects/work/kcp/code/kcp/certs/ca-cert.pem
  proxy_client_cert: certs/proxy-cert.pem
  proxy_client_key: certs/proxy-key.pem
- path: /application/services
  backend: https://localhost:6444
  backend_server_ca: /home/<username>/projects/work/kcp/code/kcp/certs/ca-cert.pem
  proxy_client_cert: certs/proxy-cert.pem
  proxy_client_key: certs/proxy-key.pem
```

## How to run it
The proxy by default listens on port 8083.
```bash
$ kcp-proxy --help
Usage:
  kcp-proxy [flags]

Flags:
      --client-ca-cert string     CA cert used to validate client certs (default "certs/ca-cert.pem")
  -h, --help                      help for kcp-proxy
      --listen-address string     Address and port for the proxy to listen on (default ":8083")
      --mapping-file string       Config file mapping paths to backends
      --server-cert-file string   The proxy's serving cert file (default "certs/server-cert.pem")
      --server-key-file string    The proxy's serving private key file (default "certs/server-key.pem")

$ go install github.com/csams/kcp-proxy
$ kcp-proxy --mapping-file=path-mapping.yaml
```