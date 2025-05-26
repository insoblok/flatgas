package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Contianer represents a container of data with a version and timestamp.
type Contianer[T any] struct {
	Version int
	Time    time.Time
	Data    T
}

// Config represents the RPC configuration.
type Config struct {
	DefaultRPC string `json:"defaultRpc"`
	RPCs       map[string]string
}

// VersionConfig creates a new Contianer with a version and timestamp, then saves it to JSON.
func VersionConfig(config Config) (string, error) {
	c := Contianer[Config]{
		Version: 1,
		Time:    time.Now(),
		Data:    config,
	}

	filename := filepath.Join("config", fmt.Sprintf("config%d.json"))
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func main() {
	// Example Config
	myConfig := Config{
		DefaultRPC: "grpc",
		RPCs: map[string]string{
			"serviceA": "grpc",
			"serviceB": "http",
		},
	}

	filename, err := VersionConfig(myConfig)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Config saved to:", filename)
}
