package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== Vault Redis Cache Test ===")
	fmt.Println()

	// Check if Redis is available
	if os.Getenv("REDIS_ADDR") == "" {
		fmt.Println("❌ REDIS_ADDR environment variable not set")
		fmt.Println("Please set REDIS_ADDR to your Redis server address")
		fmt.Println("Example: export REDIS_ADDR=localhost:6379")
		os.Exit(1)
	}

	// Check if Vault binary is available
	vaultPath := findVaultBinary()
	if vaultPath == "" {
		fmt.Println("❌ Vault binary not found")
		fmt.Println("Please ensure Vault is installed and available in PATH")
		os.Exit(1)
	}

	fmt.Printf("✅ Found Vault binary: %s\n", vaultPath)
	fmt.Printf("✅ Redis address: %s\n", os.Getenv("REDIS_ADDR"))
	fmt.Println()

	// Create temporary directory for Vault data
	tempDir, err := os.MkdirTemp("", "vault-redis-test-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	vaultDataDir := filepath.Join(tempDir, "vault-data")
	os.MkdirAll(vaultDataDir, 0o755)

	// Create Vault configuration file
	configPath := filepath.Join(tempDir, "vault.hcl")
	configContent := fmt.Sprintf(`
storage "file" {
  path = "%s"
}

listener "tcp" {
  address = "127.0.0.1:9200"
  tls_disable = 1
}

# Cache backend will be selected via environment variables
# CACHE_BACKEND=redis for Redis, or omit for LRU
`, vaultDataDir)

	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	fmt.Printf("✅ Created Vault config: %s\n", configPath)
	fmt.Println()

	// Set environment variables for Redis cache
	env := os.Environ()
	env = append(env, "CACHE_BACKEND=redis")
	env = append(env, "REDIS_ADDR="+os.Getenv("REDIS_ADDR"))
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		env = append(env, "REDIS_PASSWORD="+redisPassword)
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		env = append(env, "REDIS_DB="+redisDB)
	}

	// Start Vault server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, vaultPath, "server", "-dev", "-dev-root-token-id=dev-only-token", "-config="+configPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("🚀 Starting Vault server with Redis cache...")
	fmt.Println("Environment variables:")
	fmt.Printf("  CACHE_BACKEND=%s\n", os.Getenv("CACHE_BACKEND"))
	fmt.Printf("  REDIS_ADDR=%s\n", os.Getenv("REDIS_ADDR"))
	fmt.Printf("  REDIS_PASSWORD=%s\n", os.Getenv("REDIS_PASSWORD"))
	fmt.Printf("  REDIS_DB=%s\n", os.Getenv("REDIS_DB"))
	fmt.Println()

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Vault server: %v", err)
	}

	// Wait for Vault to start
	fmt.Println("⏳ Waiting for Vault server to start...")
	time.Sleep(5 * time.Second)

	// Test Vault API
	fmt.Println("🔍 Testing Vault API...")
	testVaultAPI()

	// Keep server running for a bit to demonstrate
	fmt.Println("⏳ Keeping Vault server running for 10 seconds to demonstrate...")
	fmt.Println("Press Ctrl+C to stop early")

	// Set up signal handling for graceful shutdown
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("🛑 Stopping Vault server...")
		cancel()
	}()

	// Wait for server to stop
	cmd.Wait()
	fmt.Println("✅ Vault server stopped")
}

func testVaultAPI() {
	// Test basic Vault API functionality
	client := &http.Client{Timeout: 10 * time.Second}

	// Test health endpoint
	resp, err := client.Get("http://127.0.0.1:8200/v1/sys/health")
	if err != nil {
		fmt.Printf("❌ Vault server not responding: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Vault server is responding")
	} else {
		fmt.Printf("⚠️  Vault server responded with status: %d\n", resp.StatusCode)
	}

	// Test initialization status
	resp, err = client.Get("http://127.0.0.1:8200/v1/sys/init")
	if err != nil {
		fmt.Printf("❌ Failed to check initialization: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Vault initialization endpoint responding")
	} else {
		fmt.Printf("⚠️  Initialization endpoint status: %d\n", resp.StatusCode)
	}
}

// Helper function to find Vault binary
func findVaultBinary() string {
	// Check common locations
	paths := []string{
		"bao",
		"./bao",
		"../bao",
		"../../bao",
		"../../../bao",
		"./bin/bao",
	}

	// Check PATH
	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return path
		}
	}

	return ""
}
