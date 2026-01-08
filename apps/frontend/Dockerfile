# Build stage
FROM node:20-alpine AS build-stage

WORKDIR /app

# Build args for Vite (needed at build time)
ARG VITE_CURRENCY=IDR
ENV VITE_CURRENCY=$VITE_CURRENCY

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage
FROM nginx:stable-alpine AS production-stage

COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
