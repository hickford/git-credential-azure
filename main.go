package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

var (
	verbose bool
	// populated by GoReleaser https://goreleaser.com/cookbooks/using-main.version
	version = "dev"
)

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if ok && version == "dev" {
		version = info.Main.Version
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "git-credential-azure %s\n", version)
	}
}

func parse(input string) map[string]string {
	lines := strings.Split(string(input), "\n")
	pairs := map[string]string{}
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) >= 2 {
			pairs[parts[0]] = parts[1]
		}
	}
	return pairs
}

func main() {
	flag.BoolVar(&verbose, "verbose", false, "log debug information to stderr")
	flag.Usage = func() {
		printVersion()
		fmt.Fprintln(os.Stderr, "usage: git credential-azure [<options>] <action>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Actions:")
		fmt.Fprintln(os.Stderr, "  get            Generate credential")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "See also https://github.com/hickford/git-credential-azure")
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(2)
	}
	switch args[0] {
	case "get":
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		pairs := parse(string(input))
		if pairs["host"] != "dev.azure.com" {
			return
		}
		printVersion()
		if verbose {
			fmt.Fprintln(os.Stderr, "input:", pairs)
		}
		result, err := authenticate()
		if err != nil {
			log.Fatalln(err)
		}
		if verbose {
			fmt.Fprintln(os.Stderr, "result:", result)
		}
		organization := strings.SplitN(pairs["path"], "/", 2)[0]
		var pt PatToken
		if organization != "" {
			pt, err = getPAT(organization, result.AccessToken)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error acquiring Personal Access Token", err)
			}
		}
		var username string
		if pairs["username"] == "" {
			// TODO: check correctness
			username = "oauth2"
		}
		output := map[string]string{}
		if username != "" {
			output["username"] = username
		}
		var password string
		var expiry time.Time
		if pt.Token != "" {
			password = pt.Token
			expiry = pt.ValidTo
		} else {
			password = result.AccessToken
			expiry = result.ExpiresOn
		}
		output["password"] = password
		if !expiry.IsZero() {
			output["password_expiry_utc"] = fmt.Sprintf("%d", expiry.UTC().Unix())
		}
		if verbose {
			fmt.Fprintln(os.Stderr, "output:", output)
		}
		for key, v := range output {
			fmt.Printf("%s=%s\n", key, v)
		}
	}
}

func authenticate() (public.AuthResult, error) {
	client, err := public.New(
		// https://github.com/git-ecosystem/git-credential-manager/blob/8c430c9484c90433ab30c25df7fc1005fe2f4ba4/src/shared/Microsoft.AzureRepos/AzureDevOpsConstants.cs#L15
		// magic https://developercommunity.visualstudio.com/t/non-interactive-aad-auth-works-for-visual-studio-a/387853
		"872cd9fa-d31f-45e0-9eab-6e460a02d1f1",
		// https://learn.microsoft.com/en-us/azure/active-directory/develop/msal-client-application-configuration#authority
		// tested with personal account
		public.WithAuthority("https://login.microsoftonline.com/organizations"))
	if err != nil {
		return public.AuthResult{}, err
	}
	// https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/manage-personal-access-tokens-via-api?view=azure-devops
	scopes := []string{"499b84ac-1321-427f-aa17-267ca6975798/.default"}
	return client.AcquireTokenInteractive(context.Background(), scopes)
}

func getPAT(organization, accessToken string) (PatToken, error) {
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/tokens/pats/create?view=azure-devops-rest-7.1&tabs=HTTP
	// sadly https://github.com/microsoft/azure-devops-go-api doesn't have this function
	url := fmt.Sprintf("https://vssps.dev.azure.com/%s/_apis/tokens/pats?api-version=7.1-preview.1", organization)
	j := map[string]any{
		"scope": "vso.code_write vso.packaging",
	}
	body, err := json.Marshal(j)
	if err != nil {
		return PatToken{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return PatToken{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return PatToken{}, err
	}
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return PatToken{}, err
	}
	if verbose {
		fmt.Fprintln(os.Stderr, string(body))
	}
	ptr := PatTokenResult{}
	err = json.Unmarshal(body, &ptr)
	if err != nil {
		return PatToken{}, err
	}
	if ptr.PatTokenError != "" && ptr.PatTokenError != "none" {
		return PatToken{}, errors.New(ptr.PatTokenError)
	}
	return ptr.PatToken, nil
}

// https://learn.microsoft.com/en-us/rest/api/azure/devops/tokens/pats/create?view=azure-devops-rest-7.1&tabs=HTTP#pattokenresult
type PatTokenResult struct {
	PatToken      PatToken `json:"patToken"`
	PatTokenError string   `json:"patTokenError"`
}

// https://learn.microsoft.com/en-us/rest/api/azure/devops/tokens/pats/create?view=azure-devops-rest-7.1&tabs=HTTP#pattoken
type PatToken struct {
	Token   string    `json:"token"`
	ValidTo time.Time `json:"validTo"`
}
