# build the application
FROM golang:1.16.15-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN which migrate

# run the binary file
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate .
COPY app.env .
COPY entrypoint.sh .

COPY db/schema ./schema

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/entrypoint.sh" ]

