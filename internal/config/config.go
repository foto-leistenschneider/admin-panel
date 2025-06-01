package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ServerAddress string
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
