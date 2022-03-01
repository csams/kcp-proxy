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
    --bind-address=0.0.0.0 \
    --secure-port=6443 \
    --auto-publish-apis \
    --run-controllers \
    --tls-cert-file=/home/<username>/projects/work/kcp/code/kcp/certs/server-cert.pem \
    --tls-private-key-file=/home/<username>/projects/work/kcp/code/kcp/certs/server-key.pem \
    --requestheader-client-ca-file=/home/<username>/projects/work/kcp/code/kcp-proxy/certs/ca-cert.pem \
    --requestheader-username-headers=X-Remote-User \
    --requestheader-group-headers=X-Remote-Group \
```

## Generate the proxy certs
The `.cnf` files in hack will need updating so the certs include your IP address
in SANs. Then run this:
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
```bash
go install github.com/csams/kcp-proxy
kcp-proxy --mapping-file=path-mapping.yaml
```