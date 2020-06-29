FROM node:14.3.0-alpine AS builder

ENV ENV_NAME dev
ENV EGG_SERVER_ENV dev
ENV NODE_ENV dev
ENV NODE_CONFIG_ENV dev

WORKDIR /build

COPY package.json .

RUN npm install

ADD . /build

RUN npm run test
RUN npx webpack

FROM kulti/static-file-server:v1.9.0

COPY --from=builder /build/index.html /web/
COPY --from=builder /build/dist/bundle.js /web/dist/
COPY --from=builder /build/css/* /web/css/
