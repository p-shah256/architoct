FROM golang:alpine

WORKDIR /app
COPY . .
RUN go build -o main .

RUN mkdir -p /app/logs && chmod 755 /app/logs

CMD ["./main"]
