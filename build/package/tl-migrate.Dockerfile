FROM migrate/migrate:v4.11.0

COPY migrations /migrations
COPY migrate-entrypoint.sh /

ENTRYPOINT [ "/migrate-entrypoint.sh" ]

CMD ["--help"]
