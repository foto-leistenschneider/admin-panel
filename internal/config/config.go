package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ServerAddress        string
	WorkosClientId       string
	WorkosApiKey         string
	WorkosCookiePassword string
	BackupDir            string
)

func init() {
	loadEnvFile()

	if portStr, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %s", portStr)
			os.Exit(1)
		}
		ServerAddress = fmt.Sprintf(":%d", port)
	}

	if workosClientId, ok := os.LookupEnv("WORKOS_CLIENT_ID"); ok {
		WorkosClientId = workosClientId
	} else {
		fmt.Fprintf(os.Stderr, "WORKOS_CLIENT_ID is not set")
		os.Exit(1)
	}

	if workosApiKey, ok := os.LookupEnv("WORKOS_API_KEY"); ok {
		WorkosApiKey = workosApiKey
	} else {
		fmt.Fprintf(os.Stderr, "WORKOS_API_KEY is not set")
		os.Exit(1)
	}

	if workosCookiePassword, ok := os.LookupEnv("WORKOS_COOKIE_PASSWORD"); ok {
		WorkosCookiePassword = workosCookiePassword
	} else {
		fmt.Fprintf(os.Stderr, "WORKOS_COOKIE_PASSWORD is not set")
		os.Exit(1)
	}

	if backupDir, ok := os.LookupEnv("BACKUP_DIR"); ok {
		BackupDir = backupDir
	} else {
		BackupDir = "./backups"
	}
	_ = os.MkdirAll(BackupDir, 0755)
}

func loadEnvFile() {
	if _, err := os.Stat(".env"); err == nil {
		if f, err := os.Open(".env"); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					fmt.Fprintf(os.Stderr, "Error parsing .env line: %v\n", line)
					os.Exit(1)
				}
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	}
}
