FROM golang:alpine
LABEL authors="markjforte"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /opnsense-proxy-api

EXPOSE 9657
CMD ["/opnsense-proxy-api"]