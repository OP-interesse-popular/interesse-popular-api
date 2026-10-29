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

## Parlamentares

`GET /api/v1/parliamentarians/right-wing` lê o Redis e retorna parlamentares em
exercício cujos partidos possuem classificação ativa de `direita` liberal ou extrema.
A API pública não acessa PostgreSQL, fontes legislativas ou a API interna: a API de
classificação monta e atualiza o cache previamente.
Aceita `page` e `limit`, conforme o `go-commons`; por exemplo:
`/api/v1/parliamentarians/right-wing?page=2&limit=20`.
Use `party` para limitar a uma ou mais siglas: 
`/api/v1/parliamentarians/right-wing/camara?party=PL,NOVO&page=1&limit=20`.
Também aceita parâmetros repetidos: `?party=PL&party=NOVO`.
Use `state` para filtrar uma ou mais UFs: `?state=SP,RJ` ou `?state=SP&state=RJ`.

Também há rotas separadas por casa legislativa:

- `GET /api/v1/parliamentarians/right-wing/camara`
- `GET /api/v1/parliamentarians/right-wing/senado`

## Logs

Além dos logs de acesso do Gin, falhas de leitura do Redis são emitidas em JSON sem
registrar dados sensíveis.

## Formato de respostas

Listas simples seguem o padrão do `go-commons`:

```json
{"items": [], "total": 0}
```

O endpoint de parlamentares usa a resposta paginada do `go-commons`:

```json
{"items": [], "page_count": 0, "page": 1}
```

Confirmações simples retornam `{"message":"success"}`. Falhas retornam
`{"status": 500, "messages": ["Internal server error"]}` sem detalhes internos.

## Testes

```bash
go test ./...
```
