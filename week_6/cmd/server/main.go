package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/flyrank/week_6/pkg/triage"
)

func loadEnv() {
	dir := "."
	for i := 0; i < 4; i++ {
		path := filepath.Join(dir, ".env")
		file, err := os.Open(path)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
						value = strings.Trim(value, "\"")
					}
					if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
						value = strings.Trim(value, "'")
					}
					os.Setenv(key, value)
				}
			}
			return
		}
		dir = filepath.Join(dir, "..")
	}
}

func main() {
	loadEnv()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY environment variable is not set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := triage.NewClient(apiKey)

	http.HandleFunc("/api/v1/triage-ticket", client.HandleTriageTicket)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Starting Ticket Triage API on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
