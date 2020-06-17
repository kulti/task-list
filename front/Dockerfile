FROM node:14.3.0-alpine AS builder

ENV ENV_NAME dev
ENV EGG_SERVER_ENV dev
ENV NODE_ENV dev
ENV NODE_CONFIG_ENV dev

WORKDIR /build

COPY package.json .

RUN npm install

ADD . /build

RUN npx webpack

FROM halverneus/static-file-server

COPY --from=builder /build/index.html /web/
COPY --from=builder /build/dist/bundle.js /web/dist/
COPY --from=builder /build/css/* /web/css/
