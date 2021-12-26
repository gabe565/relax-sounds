ARG GO_VERSION=1.17
ARG NODE_VERSION=14

FROM golang:$GO_VERSION-alpine as go-builder
WORKDIR /app

RUN apk add --no-cache gcc g++ lame-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s"


FROM node:$NODE_VERSION-alpine AS node-builder
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
CMD ["./relax-sounds", "--static=dist"]
