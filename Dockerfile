FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS api-build
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


FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-deps
WORKDIR /app

RUN corepack enable

COPY frontend/package.json frontend/pnpm-*.yaml frontend/.npmrc ./
RUN --mount=type=cache,id=pnpm,target=/root/.cache \
  pnpm install --prod --frozen-lockfile

FROM frontend-deps AS frontend-build
COPY frontend/ ./
RUN --mount=type=cache,id=pnpm,target=/root/.cache \
  pnpm run build


FROM alpine:3.22.0 AS api-only
LABEL org.opencontainers.image.source="https://github.com/gabe565/relax-sounds"
WORKDIR /app

RUN apk add --no-cache ffmpeg lame-libs tzdata

ARG USERNAME=relax-sounds
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

RUN mkdir pb_data && chown 1000:1000 pb_data

COPY --from=api-build /app/relax-sounds ./

USER $UID
CMD ["./relax-sounds", "serve", "--http=0.0.0.0:80", "--dir=/data", "--public=public"]

FROM api-only
COPY --from=frontend-build /app/dist public/
