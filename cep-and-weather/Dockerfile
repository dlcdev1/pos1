#FROM golang:1.23.4 as build
#WORKDIR /app
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun
#
#FROM scratch
#WORKDIR /app
#COPY --from=build /app/cloudrun .
#ENTRYPOINT ["./cloudrun"]

FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cep-and-weather .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/cep-and-weather .
CMD ["./cep-and-weather"]
