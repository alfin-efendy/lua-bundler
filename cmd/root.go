package cmd

import (
	"fmt"
	"os"

	"github.com/alfin-efendy/lua-bundler/internal/bundler"
	httpserver "github.com/alfin-efendy/lua-bundler/internal/http"
	"github.com/alfin-efendy/lua-bundler/internal/obfuscator"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
	// Styles using Lipgloss
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61DAFB")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
)

var rootCmd = &cobra.Command{
	Use:   "lua-bundler",
	Short: "A beautiful CLI tool for bundling Lua scripts",
	Long: lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(" Lua Script Bundler "),
		"",
		"Bundle multiple Lua files and HTTP dependencies into a single script.",
		"Supports local modules, HTTP modules via game:HttpGet(), and more!",
		"",
		infoStyle.Render("Features:"),
		"  • Bundle local Lua modules with require()",
		"  • Embed HTTP dependencies from game:HttpGet()",
		"  • Release mode to remove debug statements",
		"  • Code obfuscation support (3 levels)",
		"  • HTTP server to serve bundled output",
		"  • Beautiful terminal output with colors",
		"",
		warningStyle.Render("Example:"),
		"  lua-bundler -e main.lua -o bundle.lua --release --obfuscate 2",
		"  lua-bundler -e main.lua -o bundle.lua --serve --port 8080",
	),
	Run: func(cmd *cobra.Command, args []string) {
		entryFile, _ := cmd.Flags().GetString("entry")
		outputFile, _ := cmd.Flags().GetString("output")
		release, _ := cmd.Flags().GetBool("release")
		verbose, _ := cmd.Flags().GetBool("verbose")
		obfuscateLevel, _ := cmd.Flags().GetInt("obfuscate")
		serve, _ := cmd.Flags().GetBool("serve")
		port, _ := cmd.Flags().GetInt("port")

		if entryFile == "" {
			fmt.Println(errorStyle.Render("❌ Entry file is required"))
			os.Exit(1)
		}

		// Print header
		fmt.Println(titleStyle.Render(" Lua Script Bundler "))
		fmt.Println()
		fmt.Println(infoStyle.Render("Configuration:"))
		fmt.Printf("  Entry: %s\n", entryFile)
		fmt.Printf("  Output: %s\n", outputFile)
		if release {
			fmt.Printf("  Mode: %s\n", warningStyle.Render("Release (debug statements removed)"))
		} else {
			fmt.Printf("  Mode: %s\n", infoStyle.Render("Development"))
		}
		if obfuscateLevel > 0 {
			levelName := []string{"None", "Basic", "Medium", "Heavy"}
			if obfuscateLevel > 3 {
				obfuscateLevel = 3
			}
			fmt.Printf("  Obfuscation: %s\n", warningStyle.Render(levelName[obfuscateLevel]))
		}
		if verbose {
			fmt.Printf("  Verbose: %s\n", infoStyle.Render("Enabled"))
		}
		if serve {
			fmt.Printf("  HTTP Server: %s\n", infoStyle.Render(fmt.Sprintf("Port %d", port)))
		}
		fmt.Println()

		// Create bundler
		b, err := bundler.NewBundler(entryFile, verbose)
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to create bundler: %v", err)))
			os.Exit(1)
		}

		// Bundle
		fmt.Println(infoStyle.Render("🔄 Processing dependencies..."))
		result, err := b.Bundle(release)
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Bundling failed: %v", err)))
			os.Exit(1)
		}

		// Apply obfuscation if enabled
		if obfuscateLevel > 0 {
			fmt.Println(infoStyle.Render("🔒 Applying obfuscation..."))
			obf := obfuscator.NewObfuscator(obfuscateLevel)
			result = obf.Obfuscate(result)
		}

		// Write output
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to write output: %v", err)))
			os.Exit(1)
		}

		// Success message
		printSuccess(b, outputFile, obfuscateLevel)

		// Start HTTP server if serve flag is enabled
		if serve {
			httpserver.StartServer(outputFile, port)
		}
	},
}

func printSuccess(b *bundler.Bundler, outputFile string, obfuscateLevel int) {
	fmt.Println()
	fmt.Println(successStyle.Render("✅ Successfully bundled!"))
	fmt.Printf("%s %d\n",
		infoStyle.Render("📦 Modules embedded:"),
		len(b.GetModules()))

	if obfuscateLevel > 0 {
		fmt.Printf("%s Level %d applied\n",
			infoStyle.Render("🔒 Obfuscation:"),
			obfuscateLevel)
	}

	fmt.Printf("%s %s\n",
		successStyle.Render("📄 Output:"),
		outputFile)
}

// SetVersionInfo sets the version information from build-time variables
func SetVersionInfo(v, date, commit string) {
	version = v
	buildDate = date
	gitCommit = commit

	// Truncate commit hash for display (handle short commits)
	commitDisplay := commit
	if len(commit) > 8 {
		commitDisplay = commit[:8]
	}

	// Update root command with version
	rootCmd.Version = fmt.Sprintf("%s (built: %s, commit: %s)", version, buildDate, commitDisplay)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ %v", err)))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("entry", "e", "main.lua", "Entry point Lua file")
	rootCmd.Flags().StringP("output", "o", "bundle.lua", "Output bundled file")
	rootCmd.Flags().BoolP("release", "r", false, "Release mode: remove print and warn statements")
	rootCmd.Flags().IntP("obfuscate", "O", 0, "Obfuscation level (0=none, 1=basic, 2=medium, 3=heavy)")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolP("serve", "s", false, "Start HTTP server to serve the output file")
	rootCmd.Flags().IntP("port", "p", 8080, "Port for HTTP server (used with --serve)")
}
