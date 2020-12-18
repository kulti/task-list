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

## Deployment

### Yandex Disk

1. Creates application with app folder permissions here https://oauth.yandex.com/.
2. Use any URL as a Callback URL.
3. Make a request in browser: https://oauth.yandex.com/authorize?response_type=code&client_id=<client_id>.
4. Allow and copy code after redirect.
5. Convert code to token: `curl -X POST -u <client_id>:<password> 'https://oauth.yandex.com/token' -d 'grant_type=authorization_code&code=<code>'`.
6. Set up the TOKEN environment variable.
