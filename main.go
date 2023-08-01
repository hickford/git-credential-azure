package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"

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
		var username string
		if pairs["username"] == "" {
			username = "oauth2"
		}
		output := map[string]string{
			"password": result.AccessToken,
		}
		if username != "" {
			output["username"] = username
		}
		if !result.ExpiresOn.IsZero() {
			output["password_expiry_utc"] = fmt.Sprintf("%d", result.ExpiresOn.UTC().Unix())
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
		"872cd9fa-d31f-45e0-9eab-6e460a02d1f1",
		public.WithAuthority("https://login.microsoftonline.com/organizations"))
	if err != nil {
		return public.AuthResult{}, err
	}
	scopes := []string{"499b84ac-1321-427f-aa17-267ca6975798/.default"}
	return client.AcquireTokenInteractive(context.Background(), scopes)
}
