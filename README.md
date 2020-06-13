![CI](https://github.com/kulti/task-list/workflows/CI/badge.svg)

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
