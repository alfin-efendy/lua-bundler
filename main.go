package main

import (
	"github.com/alfin-efendy/lua-bundler/cmd"
)

// Version information (injected during build)
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.SetVersionInfo(version, buildDate, gitCommit)
	cmd.Execute()
}
