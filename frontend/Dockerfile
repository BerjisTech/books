# Build stage
FROM node:20-alpine AS build
WORKDIR /workspace
COPY clients/angular-auth ./clients/angular-auth
WORKDIR /workspace/books/frontend
COPY books/frontend/package.json ./package.json
COPY books/frontend/yarn.lock ./yarn.lock
RUN corepack enable \
 && if [ -f yarn.lock ]; then yarn install --immutable; else yarn install; fi \
 && ln -s /workspace/books/frontend/node_modules /workspace/node_modules
COPY books/frontend/ .
RUN yarn build

# Runtime stage
FROM nginx:1.27-alpine
# Angular application builder outputs browser assets under dist/books/browser
COPY --from=build /workspace/books/frontend/dist/books/browser /usr/share/nginx/html
COPY books/frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
