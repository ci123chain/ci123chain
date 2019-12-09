package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	accountRpc "github.com/tanhuiya/ci123chain/pkg/account/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/cmd/rpc"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/ibc"
	txRpc "github.com/tanhuiya/ci123chain/pkg/transfer/rest"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	distr "github.com/tanhuiya/ci123chain/pkg/distribution"
)

const (
	FlagListenAddr         = "laddr"
	FlagMaxOpenConnections = "max-open"
	FlagRPCReadTimeout     = "read-timeout"
	FlagRPCWriteTimeout    = "write-timeout"
)


func init() {
	rootCmd.AddCommand(rpcCmd)
	rpcCmd.Flags().String(FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	rpcCmd.Flags().Uint(FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	rpcCmd.Flags().Uint(FlagRPCReadTimeout, 10, "The RPC read timeout")
	rpcCmd.Flags().Uint(FlagRPCWriteTimeout, 10, "The RPC write timeout")
}

var rpcCmd = &cobra.Command{
	Use: "rest-server",
	Short: "Start rpc server",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		rs := NewRestServer()
		err := rs.Start(
			viper.GetString(FlagListenAddr),
			viper.GetInt(FlagMaxOpenConnections),
			uint(viper.GetInt(FlagRPCReadTimeout)),
			uint(viper.GetInt(FlagRPCWriteTimeout)),
			)
		return err
	},
}


// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux     *mux.Router
	CliCtx  context.Context
	listener net.Listener
}

func NewRestServer() *RestServer {
	r := mux.NewRouter()
	cliCtx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return nil
	}

	r.NotFoundHandler = Handle404()

	rpc.RegisterRoutes(cliCtx, r)
	accountRpc.RegisterRoutes(cliCtx, r)
	txRpc.RegisterTxRoutes(cliCtx, r)
	ibc.RegisterRoutes(cliCtx, r)
	distr.RegisterRoutes(cliCtx, r)
	return &RestServer{
		Mux: r,
		CliCtx: cliCtx,
	}
}

const CorePrefix = "/core"
func Handle404() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {

		nodeUri := req.RequestURI
		if strings.HasPrefix(nodeUri, CorePrefix) {
			arr := strings.SplitAfter(nodeUri, CorePrefix)
			arr = arr[1:]
			newPath := strings.Join(arr, "")
			dest := viper.GetString(helper.FlagNode)
			dest = strings.ReplaceAll(dest, "tcp", "http")

			url, _ := url.Parse(dest)
			proxy := httputil.NewSingleHostReverseProxy(url)
			req.URL.Host = url.Host
			req.URL.Path = newPath
			req.RequestURI = newPath
			// Note that ServeHttp is non blocking and uses a go routine under the hood
			proxy.ServeHTTP(w, req)
		}
	})
}

func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint) (err error) {
	util.TrapSignal(func() {
		err := rs.listener.Close()
		fmt.Println("error closing listener %v", err)
	})

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	return rpcserver.StartHTTPServer(rs.listener, rs.Mux, logger,cfg)
}

