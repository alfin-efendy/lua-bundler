package httpserver

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
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

// StartServer starts an HTTP server to serve the bundled output file
func StartServer(outputFile string, port int) {
	absPath, err := filepath.Abs(outputFile)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to get absolute path: %v", err)))
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(infoStyle.Render("🌐 Starting HTTP server..."))
	fmt.Println()

	// Get local IP addresses
	localIPs := getLocalIPs()

	// Print all access URLs
	fmt.Printf("%s http://localhost:%d/%s\n",
		successStyle.Render("🔗 Local:"),
		port,
		filepath.Base(outputFile))

	for _, ip := range localIPs {
		fmt.Printf("%s http://%s:%d/%s\n",
			successStyle.Render("🌍 Network:"),
			ip,
			port,
			filepath.Base(outputFile))
	}

	fmt.Printf("%s http://localhost:%d\n",
		infoStyle.Render("📋 Directory listing:"),
		port)
	fmt.Println()
	fmt.Println(warningStyle.Render("Press Ctrl+C to stop the server"))
	fmt.Println()

	// Create HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log request
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("[%s] %s %s %s from %s\n",
			timestamp,
			infoStyle.Render("→"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr)

		// If requesting the specific file directly
		if r.URL.Path == "/"+filepath.Base(outputFile) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.ServeFile(w, r, absPath)
			return
		}

		// If requesting root, serve directory listing
		if r.URL.Path == "/" {
			dir := filepath.Dir(absPath)
			files, err := os.ReadDir(dir)
			if err != nil {
				http.Error(w, "Unable to read directory", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<html><head><title>Lua Bundler - Output Files</title>")
			fmt.Fprintf(w, "<style>body{font-family:monospace;margin:40px;background:#1a1a1a;color:#fafafa}")
			fmt.Fprintf(w, "a{color:#61dafb;text-decoration:none;padding:5px;display:block}")
			fmt.Fprintf(w, "a:hover{background:#333;border-radius:3px}</style></head>")
			fmt.Fprintf(w, "<body><h1 style='color:#7D56F4'>📦 Lua Bundler Output Files</h1><hr><ul style='list-style:none;padding:0'>")

			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".lua" {
					fmt.Fprintf(w, "<li>📄 <a href='/%s'>%s</a></li>", file.Name(), file.Name())
				}
			}

			fmt.Fprintf(w, "</ul></body></html>")
			return
		}

		// Try to serve other files in the same directory
		requestedPath := filepath.Join(filepath.Dir(absPath), filepath.Base(r.URL.Path))
		if _, err := os.Stat(requestedPath); err == nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.ServeFile(w, r, requestedPath)
			return
		}

		http.NotFound(w, r)
	})

	// Start server on 0.0.0.0 to accept connections from any network interface
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to start server: %v", err)))
		os.Exit(1)
	}
}

// getLocalIPs returns a list of local IP addresses (excluding loopback)
func getLocalIPs() []string {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range interfaces {
		// Skip down interfaces
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Skip loopback and non-IPv4 addresses
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Only include IPv4 addresses
			ip = ip.To4()
			if ip != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips
}
