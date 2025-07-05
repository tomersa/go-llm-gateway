package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type ProviderInfo struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

type ConfigFile struct {
	VirtualKeys map[string]ProviderInfo `json:"virtual_keys"`
}

func findKeysJson() string {
	// look for keys.json in current dir, then one level up
	paths := []string{"keys.json", "../keys.json"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	fmt.Println("keys.json not found")
	os.Exit(1)
	return ""
}

func listKeys() {
	f, err := os.Open(findKeysJson())
	if err != nil {
		fmt.Println("Failed to open keys.json:", err)
		os.Exit(1)
	}
	defer f.Close()
	var cfg ConfigFile
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		fmt.Println("Failed to decode keys.json:", err)
		os.Exit(1)
	}
	fmt.Println("Available Virtual Keys:")
	for k, v := range cfg.VirtualKeys {
		fmt.Printf("- %s (provider: %s)\n", k, v.Provider)
	}
}

func runWithKey(virtualKey string) {
	jsonBody := []byte(`{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello, how are you?"}]}`)
	url := "http://localhost:8080/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+virtualKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:")
	io.Copy(os.Stdout, resp.Body)
}

func checkHealth() {
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Println("Failed to reach /health endpoint:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:")
	io.Copy(os.Stdout, resp.Body)
}

func checkMetrics() {
	resp, err := http.Get("http://localhost:8080/metrics")
	if err != nil {
		fmt.Println("Failed to reach /metrics endpoint:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:")
	io.Copy(os.Stdout, resp.Body)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Client for LLM Gateway",
	}

	var runCmd = &cobra.Command{
		Use:   "run <virtual_key>",
		Short: "Run a request with the specified virtual key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runWithKey(args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available virtual keys and their providers",
		Run: func(cmd *cobra.Command, args []string) {
			listKeys()
		},
	}

	var healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Check the health status of all configured AI providers",
		Run: func(cmd *cobra.Command, args []string) {
			checkHealth()
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(healthCmd)

	var metricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Show basic usage statistics from the server",
		Run: func(cmd *cobra.Command, args []string) {
			checkMetrics()
		},
	}
	rootCmd.AddCommand(metricsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
