FROM node:12-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
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
CMD ["nginx", "-g", "daemon off;"]
