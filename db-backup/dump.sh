#!/bin/sh

set -e

BASEDIR=/backups
TMPFILENAME=$(mktemp ${BASEDIR}/dump.XXXXXX)

trap "rm ${TMPFILENAME} &> /dev/null" EXIT

if [[ -z "${FILENAME}" ]]; then
    DATE=$(date +%d-%m-%Y)
    MONTH=$(date +%m)
    YEAR=$(date +%Y)
    BACKUPDIR=${BASEDIR}/${YEAR}/${MONTH}
    FILENAME=${BACKUPDIR}/${DATE}.dump

    mkdir -p ${BACKUPDIR}
    cd ${BACKUPDIR}
fi

LOGOUTPUT=${LOGOUTPUT:-/var/log/cron.log}

echo "Backup running to ${FILENAME}" >> ${LOGOUTPUT}

PGPASSWORD=${POSTGRES_PASSWORD} pg_dump --username=${POSTGRES_USER} --dbname=${POSTGRES_DB} --host=${POSTGRES_HOST} --data-only --file=${TMPFILENAME}

sed -i -e 's/COPY public.task_lists /DELETE FROM public.task_lists;\'$'\nCOPY public.task_lists /' \
    -e '/COPY public.schema_migrations /{N;N;d;}' \
${TMPFILENAME}

mv ${TMPFILENAME} ${FILENAME}
