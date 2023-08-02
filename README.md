git-credential-azure
====================

git-credential-azure is a Git credential helper that authenticates to [Azure Repos](https://azure.microsoft.com/en-us/products/devops/repos) (dev.azure.com). Azure Repos is part of Azure DevOps.

### Caveats

This is alpha-release software early in development:

* Untested with work and school Microsoft accounts.
* Tokens expire after 1 hour so you have to reauthenticate regularly.

A mature alternative is [Git Credential Manager](https://github.com/GitCredentialManager/git-credential-manager).

## Installation

**Download** binary from https://github.com/hickford/git-credential-azure/releases.

Then test that Git can find the application:

	git credential-azure

If you have problems, make sure that the binary is [located in the path](https://superuser.com/a/284351/62691) and [is executable](https://askubuntu.com/a/229592/18504).

### Go users

Go users can install the latest release to `~/go/bin` with:

	go install github.com/hickford/git-credential-azure@latest

## Configuration

This assumes you already have a storage helper configured such as cache or wincred.

```sh
git config --global --add credential.helper azure
```

### Unconfiguration

Run:

	git config --global --unset-all credential.helper azure

## Development

Install locally with `go install .`.

### Debugging

Use the `-verbose` flag to print more details:

```sh
git config --global --unset-all credential.helper azure
git config --global --add credential.helper "azure -verbose"
```

## See also

* [git-credential-oauth](https://github.com/hickford/git-credential-oauth): a Git credential helper that authenticates to GitHub, GitLab, BitBucket and Gerrit
* [Git Credential Manager](https://github.com/git-ecosystem/git-credential-manager)
