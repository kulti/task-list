#!/bin/sh

DSN="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_DB_HOST:5432/$POSTGRES_DB?sslmode=disable"

for i in $(seq 1 5); do
    migrate -path=/migrations -database=${DSN} $* && break
    sleep 1
done
