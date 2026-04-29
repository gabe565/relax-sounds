FROM --platform=$BUILDPLATFORM node:24-alpine AS frontend-deps
WORKDIR /app

RUN corepack enable

COPY frontend/package.json frontend/pnpm-*.yaml frontend/.npmrc ./
RUN --mount=type=cache,id=pnpm,target=/root/.cache \
  pnpm install --frozen-lockfile

COPY frontend/ ./

FROM frontend-deps AS frontend-build
RUN --mount=type=cache,id=pnpm,target=/root/.cache \
  pnpm run build


FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.9.0 AS xx
FROM --platform=$BUILDPLATFORM golang:1.26.2-alpine AS backend-build
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
COPY frontend/*.go frontend/
COPY --from=frontend-build /app/dist ./frontend/dist

ARG CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=1 xx-go build -ldflags='-w -s' -trimpath -tags grpcnotrace


FROM alpine:3.23.3
LABEL org.opencontainers.image.source="https://github.com/gabe565/relax-sounds"
WORKDIR /app

RUN apk add --no-cache ffmpeg tzdata

ARG USERNAME=relax-sounds
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

RUN mkdir /data && chown 1000:1000 /data

COPY --from=backend-build /app/relax-sounds ./

USER $UID
CMD ["/app/relax-sounds", "serve", "--http=0.0.0.0:80", "--dir=/data"]
