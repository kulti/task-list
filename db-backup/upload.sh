#!/bin/sh

set -e -o pipefail

if [[ -z "${TOKEN}" ]]; then
    echo "empty token - no backup"
    exit 0
fi

FOLDER=$1
NAME=$2
upload_file="${NAME}.tgz"

trap "rm ${upload_file} &> /dev/null" EXIT

tar -czf ${upload_file} ${FOLDER}

upload_url=$(curl -H "Authorization: OAuth ${TOKEN}" \
    "https://cloud-api.yandex.net/v1/disk/resources/upload?path=app:/${upload_file}&overwrite=true" \
    | jq --raw-output '.href')

 curl -T ${upload_file} ${upload_url}
