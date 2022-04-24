FROM golang:1.16.15-alpine3.15
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/app/main" ]

