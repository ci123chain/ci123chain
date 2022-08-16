package upgrade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/otiai10/copy"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DoUpgrade will be called after the log message has been parsed and the process has terminated.
// We can now make any changes to the underlying directory without interference and leave it
// in a state, so we can make a proper restart
func DoUpgrade(cfg *Config, bin string) error {
	// Get current node height
	height, err := GetHeight(cfg, bin)
	if err != nil {
		return fmt.Errorf("cannot get current height: %w", err)
	}

	// Get current node version
	version, err := GetVersion(cfg, bin)
	if err != nil {
		return fmt.Errorf("cannot get current version: %w", err)
	}

	// If height and version not matching, then we try to download it... maybe
	if err := DownloadBinary(cfg, height, version, bin); err != nil {
		return fmt.Errorf("cannot download binary: %w", err)
	}

	// and then set the binary again
	if err := EnsureBinary(bin); err != nil {
		return fmt.Errorf("downloaded binary doesn't check out: %w", err)
	}

	return nil
}

// DownloadBinary will grab the binary and place it in the proper directory
func DownloadBinary(cfg *Config, height uint64, version string, binPath string) error {
	url, err := GetDownloadURL(cfg, height, version)
	if err != nil {
		return err
	}

	if url != "" {
		newBinPath := cfg.NewBin()
		_ = os.Remove(newBinPath)
		err = getter.GetFile(newBinPath, url)
		if err != nil {
			return err
		}

		err = HealthCheck(newBinPath)
		if err != nil {
			return err
		}

		err = copy.Copy(newBinPath, binPath)
		if err != nil {
			return err
		}
	}

	// if it is successful, let's ensure the binary is executable
	return MarkExecutable(binPath)
}

// MarkExecutable will try to set the executable bits if not already set
// Fails if file doesn't exist or we cannot set those bits
func MarkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stating binary: %w", err)
	}
	// end early if world exec already set
	if info.Mode()&0001 == 1 {
		return nil
	}
	// now try to set all exec bits
	newMode := info.Mode().Perm() | 0111
	return os.Chmod(path, newMode)
}

// GetDownloadURL will check if there is an arch-dependent binary specified in Info
func GetDownloadURL(cfg *Config, height uint64, version string) (string, error) {
	req, err := http.NewRequest("POST", cfg.GetUpgradeUrl(height), nil)
	if err != nil {
		return "", fmt.Errorf("create get download url request err: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get can upgrade version download url err: %w", err)
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read upgrade version download url err: %w", err)
	}

	var res struct {
		State uint `json:"state"`
		Data  struct {
			Url     string `json:"url"`
			Version string `json:"version"`
		}
		Msg string `json:"msg"`
	}
	err = json.Unmarshal(out, &res)
	if err != nil {
		return "", fmt.Errorf("unmarshal get download url result err: %w", err)
	}
	if res.State != 1 {
		return "", fmt.Errorf("get can upgrade version download url state: %d, msg: %s", res.State, res.Msg)
	}
	if normalizeVersion(res.Data.Version) != version {
		fmt.Printf("current version: %s, current height: %d, current height version should be: %s, upgeade from: %s\n",
			version, height, normalizeVersion(res.Data.Version), res.Data.Url)
		return res.Data.Url, nil
	}
	return "", nil
}

// EnsureBinary ensures the file exists and is executable, or returns an error
func EnsureBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat dir %s: %w", path, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", info.Name())
	}

	// this checks if the world-executable bit is set (we cannot check owner easily)
	exec := info.Mode().Perm() & 0001
	if exec == 0 {
		return fmt.Errorf("%s is not world executable", info.Name())
	}

	return nil
}

func GetHeight(cfg *Config, bin string) (uint64, error) {
	var out bytes.Buffer
	cmd := exec.Command(bin, []string{"store-height", "--home", cfg.Home}...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	height, err := strconv.ParseUint(strings.TrimRight(out.String(), "\n"), 10, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}

func GetVersion(cfg *Config, bin string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(bin, []string{"version", "--short", "--home", cfg.Home}...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return normalizeVersion(out.String()), nil
}

func normalizeVersion(version string) string {
	return strings.Split(strings.TrimLeft(strings.ReplaceAll(version, "\n", ""), "v"), "-")[0]
}

func HealthCheck(bin string) error {
	err := MarkExecutable(bin)
	if err != nil {
		return err
	}

	cmd := exec.Command(bin)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
