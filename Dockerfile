FROM golang:1.23.1

WORKDIR /vk_bot

RUN apt-get update && apt-get install -y \
    libssl-dev \
    pkg-config \
    git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o app ./cmd/main.go

CMD ["./app"]
