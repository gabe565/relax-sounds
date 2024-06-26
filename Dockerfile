FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine AS go-dependencies
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


FROM golang:1.22.3-alpine AS go-builder
WORKDIR /app

RUN apk add --no-cache gcc g++ lame-dev

COPY --from=go-dependencies /app /app
COPY --from=go-dependencies /go /go

COPY *.go ./
COPY migrations/ migrations/
COPY internal/ internal/

ARG CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN --mount=type=cache,target=/root/.cache \
    go build -ldflags="-w -s" -trimpath -tags disable_automigrate


FROM --platform=$BUILDPLATFORM node:20-alpine AS node-builder
WORKDIR /app

COPY frontend/package.json frontend/package-lock.json frontend/.npmrc ./
ARG FONTAWESOME_NPM_AUTH_TOKEN
RUN npm ci

COPY frontend/ ./
RUN npm run build


FROM alpine:3.19
LABEL org.opencontainers.image.source="https://github.com/gabe565/relax-sounds"
WORKDIR /app

RUN apk add --no-cache lame-libs tzdata ffmpeg

ARG USERNAME=relax-sounds
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

RUN mkdir pb_data && chown 1000:1000 pb_data

COPY --from=node-builder /app/dist public/
COPY --from=go-builder /app/relax-sounds ./

USER $UID
CMD ["./relax-sounds", "serve", "--http=0.0.0.0:80", "--dir=/data", "--public=public"]
