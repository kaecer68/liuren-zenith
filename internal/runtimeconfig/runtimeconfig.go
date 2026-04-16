package runtimeconfig

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const portsFileName = ".env.ports"

func Lookup(keys ...string) (string, error) {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value, nil
		}
	}

	values, err := readPortsFile()
	if err != nil {
		return "", err
	}

	for _, key := range keys {
		if value := strings.TrimSpace(values[key]); value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("missing runtime configuration for %s", strings.Join(keys, "/"))
}

func MustLookup(keys ...string) string {
	value, err := Lookup(keys...)
	if err != nil {
		panic(err)
	}
	return value
}

func MustLookupPort(contractEnvKey, legacyEnvKey string) string {
	value, err := Lookup(contractEnvKey, legacyEnvKey)
	if err != nil {
		panic(err)
	}
	if !isDigits(value) {
		panic(fmt.Errorf("invalid runtime port for %s/%s: %q", contractEnvKey, legacyEnvKey, value))
	}
	return value
}

func readPortsFile() (map[string]string, error) {
	portsPath := filepath.Clean(portsFileName)
	file, err := os.Open(portsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("missing %s; run scripts/sync-contracts.sh first", portsPath)
		}
		return nil, fmt.Errorf("read %s: %w", portsPath, err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", portsPath, err)
	}

	return values, nil
}

func isDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return value != ""
}
