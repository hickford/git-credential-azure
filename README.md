git-credential-azure
====================

git-credential-azure is a Git credential helper that authenticates to [Azure Repos](https://azure.microsoft.com/en-us/products/devops/repos) (dev.azure.com). Azure Repos is part of Azure DevOps.

The first time you authenticate, the helper opens a browser window to Microsoft login. Subsequent authentication is non interactive.

## Installation

**Download** binary from https://github.com/hickford/git-credential-azure/releases.

Then test that Git can find the application:

	git credential-azure

If you have problems, make sure that the binary is [located in the path](https://superuser.com/a/284351/62691) and [is executable](https://askubuntu.com/a/229592/18504).

### Linux

[Several Linux distributions](https://repology.org/project/git-credential-azure/versions) include a git-credential-azure package:

[![Packaging status](https://repology.org/badge/vertical-allrepos/git-credential-azure.svg?exclude_unsupported=1&header=)](https://repology.org/project/git-credential-azure/versions)

### Go users

Go users can install the latest release to `~/go/bin` with:

	go install github.com/hickford/git-credential-azure@latest

## Configuration

This assumes you already have a storage helper configured such as cache or wincred.

```sh
git config --global --add credential.helper azure
```

### Subtleties with multiple users or organizations

If you use more than one user or organization across Azure Repos, make sure that the remote URLs include usernames. This is the default if you copied the URLs from the Azure Repos web interface.

Alternatively, you can set [credential.useHttpPath](https://git-scm.com/docs/gitcredentials#Documentation/gitcredentials.txt-useHttpPath) to store separate credentials for each repo:

```sh
git config --global credential.https://dev.azure.com.useHttpPath true
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

* [Git Credential Manager](https://github.com/git-ecosystem/git-credential-manager): a Git credential helper that authenticates to Azure Repos (and other hosts)
  * Caveats: no support for Linux arm64
* [git-credential-oauth](https://github.com/hickford/git-credential-oauth) (sister project): a Git credential helper that authenticates to GitHub, GitLab, BitBucket and other hosts using OAuth
