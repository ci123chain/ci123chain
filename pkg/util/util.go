package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strconv"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
	"syscall"
)

///const CHIANID int64 = 999

var CHAINID int64

func Setup(id int64) {
	CHAINID = id
}



func BytesToUint64(v []byte) (uint64, error) {
	u, err := strconv.Atoi(string(v))
	return uint64(u), err
}

func Uint64ToBytes(u uint64) []byte {
	return []byte(fmt.Sprint(u))
}

func TxHash(b []byte) []byte {
	return crypto.Keccak256(b)
}

func SetEnvToViper(vp *viper.Viper, key string) []string {
	v := os.Getenv(key)
	var keys []string
	pairs := strings.Split(v, ",")
	for _, pair := range pairs {
		pair := strings.TrimSpace(pair)
		if len(pair) == 0 {
			break
		}
		kv := strings.Split(pair, "=")
		vp.Set(kv[0], kv[1])
		keys = append(keys, kv[0])
	}
	return keys
}

func CheckRequiredFlag(cmd *cobra.Command, names ...string) {
	for _, name := range names {
		if err := cmd.MarkFlagRequired(name); err != nil {
			panic(err)
		}
	}
}



// TrapSignal traps SIGINT and SIGTERM and terminates the server correctly.
func TrapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}
