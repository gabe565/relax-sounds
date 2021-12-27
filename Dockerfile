ARG GO_VERSION=1.17
ARG NODE_VERSION=14

FROM --platform=$BUILDPLATFORM golang:$GO_VERSION-alpine as go-builder
WORKDIR /app

RUN apk add --no-cache gcc g++ lame-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETPLATFORM
# Set Golang build envs based on Docker platform string
RUN set -x \
    && case "$TARGETPLATFORM" in \
        'linux/amd64') export GOARCH=amd64 ;; \
        'linux/arm/v6') export GOARCH=arm GOARM=6 ;; \
        'linux/arm/v7') export GOARCH=arm GOARM=7 ;; \
        'linux/arm64') export GOARCH=arm64 ;; \
        *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;; \
    esac \
    && go build -ldflags="-w -s"


FROM --platform=$BUILDPLATFORM node:$NODE_VERSION-alpine AS node-builder
WORKDIR /app

RUN apk add --no-cache g++ make python3

COPY package.json package-lock.json .npmrc ./
ARG FONTAWESOME_NPM_AUTH_TOKEN
RUN npm ci

COPY . .
RUN npm run build


FROM alpine
WORKDIR /app
RUN apk add --no-cache lame
COPY --from=go-builder /app/relax-sounds ./
COPY --from=node-builder /app/dist dist/

ENV RELAX_SOUNDS_ADDRESS ":80"
CMD ["./relax-sounds"]
