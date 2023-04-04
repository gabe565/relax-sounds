FROM --platform=$BUILDPLATFORM golang:1.20-alpine as go-dependencies
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


FROM golang:1.20-alpine as go-builder
WORKDIR /app

RUN apk add --no-cache gcc g++ lame-dev

COPY --from=go-dependencies /app /app
COPY --from=go-dependencies /go /go

COPY *.go ./
COPY migrations/ migrations/
COPY internal/ internal/

RUN --mount=type=cache,target=/root/.cache \
    go build -ldflags="-w -s"


FROM --platform=$BUILDPLATFORM node:18-alpine AS node-builder
WORKDIR /app

COPY frontend/package.json frontend/package-lock.json frontend/.npmrc ./
ARG FONTAWESOME_NPM_AUTH_TOKEN
RUN npm ci

COPY frontend/ ./
RUN npm run build


FROM alpine
LABEL org.opencontainers.image.source="https://github.com/gabe565/relax-sounds"
WORKDIR /app

RUN apk add --no-cache lame tzdata

COPY --from=go-builder /app/relax-sounds ./
COPY --from=node-builder /app/dist public/

ARG USERNAME=relax-sounds
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

CMD ["./relax-sounds", "serve", "--http=0.0.0.0:80", "--dir=/data", "--public=public"]
