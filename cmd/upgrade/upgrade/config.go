package upgrade

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config is the information passed in to control the daemon
type Config struct {
	Home          string
	Name          string
	LogBufferSize int
	UpgradeUrl    string
}

// Root returns the root directory where all info lives
func (cfg *Config) Root() string {
	return filepath.Join(cfg.Home, "bin")
}

func (cfg *Config) GetUpgradeUrl(height uint64) string {
	return cfg.UpgradeUrl + fmt.Sprintf("/api/v1/version/getUrlByHeight?height=%d", height)
}

// CurrentBin is the path to the currently selected binary (genesis if no link is set)
// This will resolve the symlink to the underlying directory to make it easier to debug
func (cfg *Config) CurrentBin() (string, error) {
	cur := filepath.Join(cfg.Root(), cfg.Name)
	_, err := os.Lstat(cur)
	if err != nil {
		return "", err
	}
	return cur, nil
}

func (cfg *Config) NewBin() string {
	return filepath.Join(cfg.Root(), "temporary", cfg.Name)
}

// GetConfigFromEnv will read the environmental variables into a config
// and then validate it is reasonable
func GetConfigFromEnv() (*Config, error) {
	cfg := &Config{
		Home:       os.Getenv("CI_HOME"),
		Name:       os.Getenv("DAEMON_NAME"),
		UpgradeUrl: os.Getenv("UPGRADE_URL"),
	}

	logBufferSizeStr := os.Getenv("DAEMON_LOG_BUFFER_SIZE")
	if logBufferSizeStr != "" {
		logBufferSize, err := strconv.Atoi(logBufferSizeStr)
		if err != nil {
			return nil, err
		}
		cfg.LogBufferSize = logBufferSize * 1024
	} else {
		cfg.LogBufferSize = bufio.MaxScanTokenSize
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate returns an error if this config is invalid.
// it enforces Home/cosmovisor is a valid directory and exists,
// and that Name is set
func (cfg *Config) validate() error {
	if cfg.Name == "" {
		return errors.New("DAEMON_NAME is not set")
	}

	if cfg.Home == "" {
		return errors.New("CI_HOME is not set")
	}

	if cfg.UpgradeUrl == "" {
		return errors.New("UPGRADE_URL is not set")
	}
	if !strings.HasPrefix(cfg.UpgradeUrl, "http") {
		cfg.UpgradeUrl = "https://" + cfg.UpgradeUrl
	}

	if !filepath.IsAbs(cfg.Home) {
		return errors.New("DAEMON_HOME must be an absolute path")
	}

	// ensure the root directory exists
	info, err := os.Stat(cfg.Root())
	if err != nil {
		return fmt.Errorf("cannot stat home dir: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", info.Name())
	}

	return nil
}
