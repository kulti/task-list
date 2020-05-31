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
