package util

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	ProofQueryPrefix = "s/k:accounts/"

	ShardDefaultProto = "tcp://"
	ShardDefaultPort = ":26657"


	//enviromnet
	CICHAINID = "CICHAINID"

	HistorySuffix = "HistoryInfo"
)

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

///account history.
type HeightUpdate struct {
	Shard   string       `json:"shard"`
	Coins   types.Coins  `json:"coins"`
}

type Heights []int64

func (h Heights) Len() int { return len(h) }

func (h Heights) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h Heights) Less(i, j int) bool { return h[i] < h[j] }

func (h Heights) Search(i int64) int64 {
	return search(h, i)
}

func search(h Heights, i int64) int64 {
	sort.Sort(h)
	if len(h) ==1  {
		if h[len(h)-1] <= i {
			return h[len(h)-1]
		}
	}
	if h[0] > i {
		return -2
	}
	if h[len(h)-1] <= i {
		return h[len(h)-1]
	}else {
		if h[len(h)/2 - 1] == i {
			return i
		}else if h[len(h)/2 -1] < i {
			if h[len(h)/2] > i {
				return h[len(h)/2-1]
			}else {
				return search(h[len(h)/2:], i)
			}
		}else {
			return search(h[:len(h)/2-1], i)
		}
	}
}

type HeightsUpdate struct {
	Height int64   `json:"height"`
	Shard  string  `json:"shard"`
}

type HeightsUpdates []HeightsUpdate

func (h HeightsUpdates) Len() int { return len(h) }

func (h HeightsUpdates) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeightsUpdates) Less(i, j int) bool { return h[i].Height< h[j].Height }

func (h HeightsUpdates) Search(i int64) HeightsUpdate {
	return searchHeight(h, i)
}

func searchHeight(h HeightsUpdates, i int64) HeightsUpdate {
	sort.Sort(h)
	if len(h) ==1  {
		if h[len(h)-1].Height <= i {
			return h[len(h)-1]
		}
	}
	if h[0].Height > i {
		return HeightsUpdate{
			Height: -1,
			Shard:  "",
		}
	}
	if h[len(h)-1].Height <= i {
		return h[len(h)-1]
	}else {
		if h[len(h)/2 - 1].Height == i {
			return h[len(h)/2 -1]
		}else if h[len(h)/2 -1].Height < i {
			if h[len(h)/2].Height > i {
				return h[len(h)/2-1]
			}else {
				return searchHeight(h[len(h)/2:], i)
			}
		}else {
			return searchHeight(h[:len(h)/2-1], i)
		}
	}
}


//type HistoryAccount struct {
//	Shard   string       `json:"shard"`
//	Account  []byte  `json:"account"`
//	Proof   *tcrypto.ProofOps `json:"proof"`
//}
//
//func NewHistoryAccount(s string, acc []byte, p *tcrypto.ProofOps) HistoryAccount{
//	return HistoryAccount{
//		Shard: s,
//		Account: acc,
//		Proof: p,
//	}
//}