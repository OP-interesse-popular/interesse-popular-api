# interesse-popular-api

Base da API do Interesse Popular, construída em Go com Gin.

## Executar localmente

```bash
cp .env.example .env
go run .
```

A API fica disponível em `http://localhost:3000`.

## Healthcheck

```bash
curl -i http://localhost:3000/api/v1/healthcheck
```

O endpoint responde `200 OK` sem corpo.

## Testes

```bash
go test ./...
```
