FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /task-api .

FROM alpine:3.22
WORKDIR /app
COPY --from=build /task-api ./task-api

EXPOSE 8080
CMD ["./task-api"]
