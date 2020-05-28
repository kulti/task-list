#!/bin/sh

migrate -path=/migrations -database=postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_DB_HOST:5432/$POSTGRES_DB?sslmode=disable $*
