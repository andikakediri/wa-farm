FROM golang:1.22-alpine

RUN apk add --no-cache git curl bash

WORKDIR /workspaces/wa-farm

EXPOSE 8080

CMD ["sleep", "infinity"]
