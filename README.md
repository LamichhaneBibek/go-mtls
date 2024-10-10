# Gerando certificados autoassinados sem CA

## Gerar chave privada

```bash
openssl genrsa -out server.key 4096
```

## gerar a requisição de certificado (csr)

```bash
openssl req -new -key server.key -out server.csr
```

## gerar o certificado

```bash
openssl x509 -req -in server.csr -out server.crt -key server.key -days 365 -sha256
```

# Gerando certificados com CA autoassinada

## Gerar chave privada da CA

```bash
openssl genrsa -out ca.key 4096
```

## gerar o certificado da CA

```bash
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt 
```

## Gerar chave privada

```bash
openssl genrsa -out server.key 4096
```

## gerar a requisição de certificado (csr)

```bash
openssl req -new -key server.key -out server.csr
```

## gerar o certificado com assinatura da CA

```bash
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 
```

# Converter certificado para formato p12

```bash
openssl pkcs12 -export -out client.p12 -inkey client.key -in client.crt
```