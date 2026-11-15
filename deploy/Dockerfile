FROM node:22-alpine AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --silent
COPY web/ ./
RUN npm run build

FROM golang:1.26-alpine AS server
WORKDIR /src/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY . /src
WORKDIR /src
RUN M=server; [ -f plugins/go.mod ] && M=plugins; cd "$M" && \
    CGO_ENABLED=0 GOWORK=off go build -trimpath -ldflags="-s -w" -o /out/logware ./cmd/logware

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 logware
WORKDIR /app
COPY --from=server /out/logware /app/logware
COPY --from=web /src/web/dist /app/web
ENV WEB_DIR=/app/web ADDR=:8080
USER logware
EXPOSE 8080
ENTRYPOINT ["/app/logware"]
