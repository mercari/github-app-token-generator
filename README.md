#  GitHub App Token Generator

## Description

> A simple github action written in go to retrieve an installation access token for an app installed into an organization.
## Usage

```yaml
name: Checkout repos
on: push
jobs:
  checkout:
    runs-on: ubuntu-latest
    steps:
    - uses: mercari/github-app-token-generator@v1
      id: get-token
      with:
        app-private-key: ${{ secrets.APP_PRIVATE_KEY }}
        app-id: ${{ secrets.APP_ID }}
        app-installation-id: ${{ secrets.APP_INSTALLATION_ID}}
        
    - name: Check out an other repo
      uses: actions/checkout@v2
      with:
        repository: owner/repo
        token: ${{ steps.get-token.outputs.token }}
```

## Inputs

| Input                 | Description                                     | Required? | Type     |
|-----------------------|-------------------------------------------------|-----------|----------|
| `app-id`              | GitHub App ID                                   | ✅         | `number` |
| `app-installation-id` | ID of the app installation to your organization | ✅         | `number` |
| `app-private-key`     | Private key of your GitHub App (Base64 encoded) | ✅         | `string` |

## Outputs

| Output  | Description                | Type     |
|---------|----------------------------|----------|
| `token` | Generated short-live token | `string` |

## Contributions

Please read the [CLA](https://www.mercari.com/cla/) carefully before submitting your contribution to Mercari. Under any circumstances, by submitting your contribution, you are deemed to accept and agree to be bound by the terms and conditions of the CLA.

## License

Copyright 2022 Mercari, Inc.

Licensed under the MIT License.

## Credits

This was originally developed by [mlioo/go-github-app-token-generator](https://github.com/mlioo/go-github-app-token-generator).
