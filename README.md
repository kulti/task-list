![CI](https://github.com/kulti/task-list/workflows/CI/badge.svg)
[![Coverage](https://coveralls.io/repos/github/kulti/task-list/badge.svg?branch=master)](https://coveralls.io/github/kulti/task-list?branch=master)

## Development

### Vscode configuration

.vscode/settings.json should contain
```
{
    "go.testTags": "integration",
    "go.buildFlags": [
        "-mod=vendor"
    ]
}
```
