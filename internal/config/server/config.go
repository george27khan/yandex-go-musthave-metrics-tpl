package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host            string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	DSN             string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	addr := flag.String("a", "localhost:8080", "Address to listen server on")
	storeInterval := flag.Duration("i", 15, "Interval to store metrics")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "Path to file storage path")
	restore := flag.Bool("r", true, "Restore metrics from a file")
	DSN := flag.String("d", "", "Dns name to use") //postgres://postgres:postgres@127.0.0.1:5432/metrics?sslmode=disable
	flag.Parse()

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		*addr = val
	}
	if val, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("STORE_INTERVAL error: %w", err)
		}
		*storeInterval = time.Duration(i)
	}
	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		*fileStoragePath = val
	}
	if val, ok := os.LookupEnv("RESTORE"); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("RESTORE error: %w", err)
		}
		*restore = b
	}
	if val, ok := os.LookupEnv("DATABASE_DSN"); ok {
		*DSN = val
	}

	cfg.Host = *addr
	cfg.StoreInterval = *storeInterval * time.Second
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
	cfg.DSN = *DSN
	return cfg, nil
}
