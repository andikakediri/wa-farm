FROM golang:1.22-alpine

RUN apk add --no-cache git curl bash sqlite-dev gcc musl-dev

WORKDIR /workspaces/wa-farm

EXPOSE 8080

CMD ["sleep", "infinity"]
