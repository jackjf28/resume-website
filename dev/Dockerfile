FROM golang:1.25-bookworm
WORKDIR /app

RUN apt-get update && \
    apt-get install -y curl nodejs npm && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

COPY package*.json ./
RUN npm install

# dependency caching
COPY go.mod go.sum ./
RUN go mod download

COPY dev/scripts/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
