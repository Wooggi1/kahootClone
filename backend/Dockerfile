# Use uma imagem base menor, mas suficiente para buildar o Go
FROM golang:1.23 AS builder

# Configurações de ambiente para o build Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Criação do diretório de trabalho
WORKDIR /app

# Copiar os arquivos de dependências Go
COPY go.mod go.sum ./

# Fazer o download das dependências do Go
RUN go mod download

# Copiar o código-fonte da aplicação
COPY . .

# Build da aplicação
RUN go build -o main ./cmd/quiz

# Segunda etapa: usar uma imagem leve para o contêiner final
FROM alpine:latest

# Instalar bibliotecas mínimas necessárias (opcional, para compatibilidade)
RUN apk add --no-cache ca-certificates

# Configuração do diretório de trabalho
WORKDIR /root/

# Copiar o binário compilado do estágio de build
COPY --from=builder /app/main .

# Expor a porta necessária
EXPOSE 3000

# Adicionar comandos para depuração de variáveis de ambiente
RUN echo "Preparing to run application" && \
    echo "Checking Environment Variables..." && \
    echo "DB_STRING: $DB_STRING" && \
    echo "JWT_SECRET: $JWT_SECRET"

# Comando para executar a aplicação
CMD ["sh", "-c", "echo DB_STRING=$DB_STRING && echo JWT_SECRET=$JWT_SECRET && ./main"]
