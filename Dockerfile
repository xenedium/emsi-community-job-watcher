FROM golang:alpine AS builder
LABEL authors="xenedium"

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/src/app ./...

FROM golang:alpine AS final
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/emsi-community-job-watcher .

CMD ["./emsi-community-job-watcher"]