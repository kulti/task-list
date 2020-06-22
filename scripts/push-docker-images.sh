#!/bin/bash

set -e

git_tag=$(git describe --tags 2> /dev/null)
if [[ ${git_tag} != "" ]]; then
    echo -n ${DOCKER_HUB_TOKEN} | docker login --username kulti --password-stdin

    services="tl-proxy tl-front tl-server tl-migrate"
    for s in ${services}; do
        image_s="kulti/${s}"
        image="${image_s}:${git_tag}"
        echo -en "${image}\t tagging ... "
        docker tag ${image_s} ${image}
        echo -n "pushing ... "
        docker push ${image}
        echo "done"
    done
fi
