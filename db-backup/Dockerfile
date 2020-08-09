FROM postgres:12-alpine

RUN apk update && apk add --no-cache jq curl

COPY crontab /etc/crontabs/root
COPY dump.sh /
COPY upload.sh /

CMD ["crond", "-f", "-d", "8"]
