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

As rotas de parlamentares leem o Redis. A API pública não acessa PostgreSQL, fontes
legislativas ou a API interna: a API de classificação monta e atualiza o cache
previamente. Deputados aceitam `page` e `limit`; por exemplo:
`/api/v1/parliamentarians/deputies?page=2&limit=20`.
Use `party` para limitar a uma ou mais siglas: 
`/api/v1/parliamentarians/deputies?party=PL,NOVO&page=1&limit=20`.
Também aceita parâmetros repetidos: `?party=PL&party=NOVO`.
Use `state` para filtrar uma ou mais UFs: `?state=SP,RJ` ou `?state=SP&state=RJ`.
Use `direction` para `esquerda` ou `direita`, e `ideology` para uma ou mais
ideologias, por exemplo: `?direction=direita&ideology=direita_liberal,extrema_direita`.

Também há rotas separadas por casa legislativa:

- `GET /api/v1/parliamentarians/deputies`
- `GET /api/v1/parliamentarians/senators`

Deputados retornam uma resposta paginada. Senadores retornam todos os itens em uma
única lista (`items` e `total`), pois a casa tem quantidade fixa e pequena de vagas.

## Políticos

`GET /api/v1/politicians` lista uma pessoa por item e inclui todos os seus mandatos
federais conhecidos. A rota usa `opip:politicians:v1`, a projeção produzida pela API
interna, e é paginada (`page`, `limit`; padrão 30, máximo 100).

Filtros disponíveis: `q` (nome), `party`, `state`, `house` (`camara` ou `senado`),
`office` (`deputado_federal` ou `senador`), `status` (`active`, `elected_pending`
ou `former`), `electionYear`, `direction`, `ideology` e `inOffice` (`true` ou
`false`). Os filtros de mandato são avaliados juntos no mesmo mandato; valores
múltiplos podem ser passados separados por vírgula ou repetindo o parâmetro. Direção
e ideologia são a classificação editorial vigente do partido associada ao mandato.

Ordenação: `sort=name|startsAt|endsAt` e `order=asc|desc`. Para datas, é usada a
data mais recente conhecida entre os mandatos da pessoa.

Exemplos:

```text
/api/v1/politicians?party=PL&state=SP&inOffice=true&page=1&limit=20
/api/v1/politicians?status=elected_pending&house=camara
/api/v1/politicians?q=Jones
```

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
