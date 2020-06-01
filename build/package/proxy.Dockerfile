FROM nginx:1.17.10-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
