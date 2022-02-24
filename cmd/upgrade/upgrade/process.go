package upgrade

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// LaunchProcess runs a subprocess and returns when the subprocess exits,
// either when it dies, or *after* a successful upgrade.
func LaunchProcess(cfg *Config, args []string, stdout, stderr io.Writer) error {
	// Get current node binary
	bin, err := cfg.CurrentBin()
	if err != nil {
		return fmt.Errorf("error creating symlink to genesis: %w", err)
	}

	if err := EnsureBinary(bin); err != nil {
		return fmt.Errorf("current binary invalid: %w", err)
	}

	if err := DoUpgrade(cfg, bin); err != nil {
		return fmt.Errorf("upgrade binary failed: %w", err)
	}

	cmd := exec.Command(bin, args...)
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	scanOut := bufio.NewScanner(io.TeeReader(outpipe, stdout))
	scanErr := bufio.NewScanner(io.TeeReader(errpipe, stderr))
	// set scanner's buffer size to cfg.LogBufferSize, and ensure larger than bufio.MaxScanTokenSize otherwise fallback to bufio.MaxScanTokenSize
	var maxCapacity int
	if cfg.LogBufferSize < bufio.MaxScanTokenSize {
		maxCapacity = bufio.MaxScanTokenSize
	} else {
		maxCapacity = cfg.LogBufferSize
	}
	bufOut := make([]byte, maxCapacity)
	bufErr := make([]byte, maxCapacity)
	scanOut.Buffer(bufOut, maxCapacity)
	scanErr.Buffer(bufErr, maxCapacity)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launching process %s %s: %w", bin, strings.Join(args, " "), err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if err := cmd.Process.Signal(sig); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for scanOut.Scan() {
			fmt.Print(scanOut.Text())
		}
	}()

	go func() {
		for scanErr.Scan() {
			fmt.Print(scanErr.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
