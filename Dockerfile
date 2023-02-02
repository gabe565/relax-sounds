FROM golang:1.20-alpine as go-builder
WORKDIR /app

RUN apk add --no-cache gcc g++ lame-dev

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
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
COPY --from=node-builder /app/dist frontend/

ARG USERNAME=relax-sounds
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

COPY data-default /data

ENV RELAX_SOUNDS_ADDRESS ":80"
ENV RELAX_SOUNDS_DATA "/data"
ENV RELAX_SOUNDS_FRONTEND "frontend"
CMD ["./relax-sounds"]
