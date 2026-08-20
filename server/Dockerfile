FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/chat-server ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/chat-server /chat-server

EXPOSE 3000
ENTRYPOINT ["/chat-server"]
