FROM golang:1.18.5-alpine3.16 as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download && go build -o ./todo-service ./cmd

FROM alpine:3.16

EXPOSE 8080
WORKDIR /apps
COPY --from=build /build/todo-service .
COPY ./.env . 

ENTRYPOINT [ "/apps/todo-service" ]
