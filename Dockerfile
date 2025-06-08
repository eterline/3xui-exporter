FROM golang:1.24.3 AS gobuilder

WORKDIR /app
COPY . .

RUN \
    go mod tidy && \ 
    CGO_ENABLED=0  \ 
    GOOS=linux     \ 
    GOARCH=amd64   \ 
    go build -o app -v ./cmd/3xui-exporter/...

FROM alpine:latest

WORKDIR /app
COPY --from=gobuilder /app/app .

RUN \
    chmod +x /app/app && \
    mkdir logs

CMD [ "./app", "--listen", "localhost:9100" ]