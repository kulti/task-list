# /bin/bash
set -o pipefail

IGNORE_PACKAGES="\
    github.com/kulti/task-list$
    github.com/kulti/task-list/cmd$
    github.com/kulti/task-list/internal/apitest$
    github.com/kulti/task-list/internal/generated
"

PACKAGES_FILTER=$(echo ${IGNORE_PACKAGES} | sed -e 's/ /|/g')

PACKAGES=$(go list -f '{{.Name}} {{.Dir}} {{.ImportPath}}' ./... | grep -v -E "${PACKAGES_FILTER}")

ALL_PACKAGES=""
IFS_BACKUP=${IFS}
IFS=$'\n'
for p in ${PACKAGES}; do
    name=$(echo $p | cut -f1 -d ' ')
    dir=$(echo $p | cut -f2 -d ' ')
    pkg=$(echo $p | cut -f3 -d ' ')

    ALL_PACKAGES="${ALL_PACKAGES} ${pkg}"
    echo "package ${name}_test" > ${dir}/empty_test.go
done

IFS=${IFS_BACKUP}
go test -v -mod=vendor -cover -covermode=atomic -coverprofile=coverage.txt ${ALL_PACKAGES} | sed -e '/testing: warning: no tests to run/{N;N;d;}'
