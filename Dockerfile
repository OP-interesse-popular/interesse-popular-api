ARG GO_VERSION=1.25.5

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/interesse-popular-api .

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=build /bin/interesse-popular-api /app/interesse-popular-api

ENV APP_ENV=production
ENV PORT=3000

EXPOSE 3000
USER app

CMD ["/app/interesse-popular-api"]
