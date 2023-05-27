FROM golang:1.19

ENV GO111MODULE=on
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main main.go
CMD ["/app/main"]