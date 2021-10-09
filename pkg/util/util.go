package util

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

///const CHIANID int64 = 999
const (
	IDG_APPID = "IDG_APPID"
	DefaultTCP = "tcp://"
	DefaultHTTP = "http://"
	DefaultHTTPS = "https://"
	DefaultWS = "ws"
	DefaultWSS = "wss"
)

var CHAINID int64
var IteratorLimit int

func SchemaPrefix() string {
	prefix := DefaultHTTP
	if os.Getenv(IDG_APPID) != "" {
		prefix = DefaultHTTPS
	}
	return prefix
}

func Setup(id int64) {
	CHAINID = id
}

func SetLimit(limit int) {
	IteratorLimit = limit
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

func GetLocalAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

///account history.
//type HeightUpdate struct {
//	Shard   string       `json:"shard"`
//	Coins   types.Coins  `json:"coins"`
//}
//
//type Heights []int64
//
//func (h Heights) Len() int { return len(h) }
//
//func (h Heights) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
//
//func (h Heights) Less(i, j int) bool { return h[i] < h[j] }
//
//func (h Heights) Search(i int64) int64 {
//	return search(h, i)
//}
//
//func search(h Heights, i int64) int64 {
//	sort.Sort(h)
//	if len(h) ==1  {
//		if h[len(h)-1] <= i {
//			return h[len(h)-1]
//		}
//	}
//	if h[0] > i {
//		return -2
//	}
//	if h[len(h)-1] <= i {
//		return h[len(h)-1]
//	}else {
//		if h[len(h)/2 - 1] == i {
//			return i
//		}else if h[len(h)/2 -1] < i {
//			if h[len(h)/2] > i {
//				return h[len(h)/2-1]
//			}else {
//				return search(h[len(h)/2:], i)
//			}
//		}else {
//			return search(h[:len(h)/2-1], i)
//		}
//	}
//}