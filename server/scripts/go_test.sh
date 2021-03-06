# /bin/bash
set -o pipefail

IGNORE_PACKAGES="\
    github.com/kulti/task-list/server\s
    github.com/kulti/task-list/server/cmd\s
    github.com/kulti/task-list/server/internal/apitest\s
    github.com/kulti/task-list/server/internal/generated
    github.com/kulti/task-list/server/internal/storages
    github.com/kulti/task-list/server/internal/services/calservice\s
"

PACKAGES_FILTER=$(echo ${IGNORE_PACKAGES} | sed -e 's/ /|/g')

PACKAGES=$(go list -f '{{.Name}} {{.Dir}} {{.ImportPath}} {{.TestGoFiles}}{{.XTestGoFiles}}' ./... | grep -v -E "${PACKAGES_FILTER}")

ALL_PACKAGES=""
IFS_BACKUP=${IFS}
IFS=$'\n'
for p in ${PACKAGES}; do
    name=$(echo $p | cut -f1 -d ' ')
    dir=$(echo $p | cut -f2 -d ' ')
    pkg=$(echo $p | cut -f3 -d ' ')
    tests=$(echo $p| cut -f4 -d ' ')

    ALL_PACKAGES="${ALL_PACKAGES} ${pkg}"
    if [[ ${tests} == '[][]' ]]; then
        echo "package ${name}_test" > ${dir}/empty_test.go
    fi
done

IFS=${IFS_BACKUP}
go test -v -mod=vendor -cover -covermode=atomic -coverprofile=coverage.txt ${ALL_PACKAGES} | sed -e '/testing: warning: no tests to run/{N;N;d;}'

echo "${ALL_PACKAGES}"
