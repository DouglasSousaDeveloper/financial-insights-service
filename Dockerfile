# Multi-stage build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependências necessárias do Alpine
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Copiar arquivos de dependência primeiro para aproveitar cache lógico do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo o código fonte
COPY . .

# Fazer a compilação garantindo um binário puramente estático (-tags netgo, CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/api

# Stage final "scratch/alpine" para imagem levíssima
FROM alpine:3.19

# Adiciona certificados SSL (necessário para comunicar com a API da OpenAI por HTTPS)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Puxa o binário do step builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
