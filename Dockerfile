FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine AS go-builder
WORKDIR /app

RUN apk add --no-cache clang

COPY --from=xx / /

COPY go.mod go.sum ./
RUN go mod download

ARG TARGETPLATFORM
RUN xx-apk add --no-cache gcc g++ lame-dev

COPY *.go ./
COPY migrations/ migrations/
COPY internal/ internal/

ARG CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=1 xx-go build -ldflags='-w -s' -trimpath -tags disable_automigrate,grpcnotrace


FROM --platform=$BUILDPLATFORM node:22-alpine AS node-builder
WORKDIR /app

COPY frontend/package.json frontend/package-lock.json frontend/.npmrc ./
RUN npm ci

COPY frontend/ ./
RUN npm run build


FROM alpine:3.21.2
LABEL org.opencontainers.image.source="https://github.com/gabe565/relax-sounds"
WORKDIR /app

RUN apk add --no-cache ffmpeg lame-libs tzdata

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
