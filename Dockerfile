FROM golang:1.22-alpine

RUN apk add --no-cache git curl screen sqlite-dev gcc musl-dev bash

WORKDIR /workspaces/wa-farm

# Pre-clone whatsmeow for faster setup
RUN git clone https://github.com/tulir/whatsmeow /workspaces/wa-farm/whatsmeow

EXPOSE 8080

CMD ["sleep", "infinity"]
