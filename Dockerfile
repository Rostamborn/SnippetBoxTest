FROM golang:1.20.0
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o main ./cmd/web/*
EXPOSE 8000
ENTRYPOINT ["./main"]
