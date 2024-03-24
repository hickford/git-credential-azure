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

### How it works

Git is cleverly designed to [support multiple credential helpers](https://git-scm.com/docs/gitcredentials#_custom_helpers). To fill credentials, Git calls each helper in turn until it has the information it needs. git-credential-azure is a read-only credential-generating helper, designed to be configured in combination with a storage helper.

To configure together with [git-credential-cache](https://git-scm.com/docs/git-credential-cache):

```sh
git config --global --unset-all credential.helper
git config --global --add credential.helper "cache --timeout 7200" # two hours
git config --global --add credential.helper azure
```

You may choose a different storage helper such as `osxkeychain`, `wincred` or `libsecret`, but git-credential-azure must be configured last. This ensures Git checks for *stored* credentials before generating *new* credentials.

**Windows users** must use storage helper `wincred` because [git-credential-cache isn't available on Windows](https://github.com/git-for-windows/git/issues/3892).

### Manual config

Edit your [global git config](https://git-scm.com/docs/git-config#FILES) `~/.gitconfig` to include the following lines:

```ini
[credential]
	helper = cache --timeout 7200	# two hours
	helper = azure
```

### Browserless systems

On systems without a web browser, set the `-device` flag to authenticate on another device using [OAuth device flow](https://www.rfc-editor.org/rfc/rfc8628).

```ini
[credential]
	helper = cache --timeout 7200	# two hours
	helper = azure -device
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
