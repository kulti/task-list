#!/bin/bash

set -e

git_tag=$(git describe --tags 2> /dev/null)
if [[ ${git_tag} != "" ]]; then
    services="tl-proxy tl-front tl-server tl-migrate"
    for s in ${services}; do
        image_s="kulti/${s}"
        image="${image_s}:${git_tag}"
        echo -en "${image}\t tagging ... "
        docker tag ${image_s} ${image}
        echo -n "pushing ... "
        docker push ${image}
    done
fi
