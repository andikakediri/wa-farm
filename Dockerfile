FROM golang:1.22-alpine

RUN apk add --no-cache git curl screen sqlite-dev gcc musl-dev

WORKDIR /app

# Clone & build whatsmeow
RUN git clone https://github.com/tulir/whatsmeow && \
    cd whatsmeow && \
    go build -o /app/wabot ./example/

# Copy server files
COPY server.go .
COPY index.html .
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go build -o /app/server server.go

EXPOSE 8080

CMD ["/app/server"]
