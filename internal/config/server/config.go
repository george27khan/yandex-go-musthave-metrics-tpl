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
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	addr := flag.String("a", "localhost:8080", "Address to listen server on")
	storeInterval := flag.Duration("i", 15, "Interval to store metrics")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "Path to file storage path")
	restore := flag.Bool("r", true, "Restore metrics from a file")
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

	cfg.Host = *addr
	cfg.StoreInterval = *storeInterval * time.Second
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
	return cfg, nil
}
