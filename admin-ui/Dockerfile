FROM node:18-alpine AS commons-builder
WORKDIR /app
RUN apk add git
RUN git init && \
    git remote add -f origin https://github.com/Coders-Brothers/backend.git && \
    git config core.sparseCheckout true && \
    echo 'commons' >> .git/info/sparse-checkout && \
    git pull origin production
WORKDIR /app/commons/ts
RUN npm install
RUN npm run build

FROM node:18-alpine AS builder
ARG CI_VERSION
ENV NEXT_PUBLIC_PRODUCT_VERSION=$CI_VERSION
ENV NODE_ENV=production
COPY --from=commons-builder /app/commons/ts/ /app/commons/ts
ADD . /app/
WORKDIR /app
RUN npm install
RUN npm install --save ./commons/ts
RUN npm run build

FROM gcr.io/distroless/nodejs18-debian11
ENV NODE_ENV=production
ENV PORT=3000
WORKDIR /app
COPY --from=builder /app/public ./public
COPY --from=builder --chown=65532:65532 /app/build/standalone ./
COPY --from=builder --chown=65532:65532 /app/build/static ./build/static
EXPOSE 3000
USER 65532:65532
CMD ["server.js"]