#!/bin/bash

set -e

git_tag=$(git describe --tags 2> /dev/null)
if [[ ${git_tag} != "" ]]; then
    major=$(echo ${git_tag} | cut -d. -f1)
    minor=$(echo ${git_tag} | cut -d. -f2)
    patch=$(echo ${git_tag} | cut -d. -f3)

    tags="${major}"
    [[ ! -z ${minor} ]] && tags="${tags} ${major}.${minor}"
    [[ ! -z ${patch} ]] && tags="${tags} ${major}.${minor}.${patch}"

    services="tl-proxy tl-front tl-server tl-migrate tl-db-backup"
    for t in ${tags}; do
        for s in ${services}; do
            image_s="kulti/${s}"
            image="${image_s}:${t}"
            echo -en "${image}\t tagging ... "
            docker tag ${image_s} ${image}
            echo -n "pushing ... "
            docker push ${image}
        done
    done
fi
