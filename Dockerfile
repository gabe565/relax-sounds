FROM node:12-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json .npmrc ./
ARG FONTAWESOME_NPM_AUTH_TOKEN
RUN set -x \
    && apk add --no-cache --virtual .build-deps \
        g++ \
        make \
        python3 \
    && npm ci \
    && apk del .build-deps

COPY . .
RUN npm run build

# Final Image
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

CMD ["nginx", "-g", "daemon off;"]
