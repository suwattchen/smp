FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /bin/core ./cmd/core

FROM alpine:3.19
RUN adduser -D appuser
USER appuser
COPY --from=build /bin/core /bin/core
EXPOSE 8080
ENTRYPOINT ["/bin/core"]
