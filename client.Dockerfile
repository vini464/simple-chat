# syntax=docker/dockerfile:1.7-labs
# primeiro passo: buildar a aplicação
FROM golang:alpine AS builder

# Cria a pasta de trabalho
WORKDIR /app
# copia os arquivos para fazer o build
COPY --exclude=server/ .. . 
# builda o programa
# CGO_ENABLED=0 -> cria o binário sem dependências externas, possibilita a execução a partir da imagem scratch
# GOOS=linux -> garante que o binario foi buildado para linux
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main client/client.go


# criação da imagem:
FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
# executa o processo
CMD ["./main"]

