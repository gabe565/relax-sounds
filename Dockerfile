    FROM node:12-alpine AS builder

WORKDIR /app
RUN apk add --no-cache g++ make python3

COPY package.json package-lock.json .npmrc ./
ARG FONTAWESOME_NPM_AUTH_TOKEN
RUN npm ci

COPY . .
RUN npm run build

# Final Image
FROM nginx:stable-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
