![CI](https://github.com/kulti/task-list/workflows/CI/badge.svg)
[![Coverage](https://coveralls.io/repos/github/kulti/task-list/badge.svg?branch=master)](https://coveralls.io/github/kulti/task-list?branch=master)
[![code style: prettier](https://img.shields.io/badge/code_style-prettier-ff69b4.svg?style=flat-square)](https://github.com/prettier/prettier)

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
