FROM postgres:12-alpine

COPY crontab /etc/crontabs/root
COPY dump.sh /

CMD ["crond", "-f", "-d", "8"]
