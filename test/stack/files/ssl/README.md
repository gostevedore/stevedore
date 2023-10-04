# README

## Generate wildcard self-signed certificate

### Generate private key and csr

```sh
openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config stevedore.test.cnf
```

### Generate certificate

```sh
openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile stevedore.test.cnf
```
