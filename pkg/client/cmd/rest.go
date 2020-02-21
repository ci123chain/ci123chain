package cmd

import (
	"encoding/json"
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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	distr "github.com/tanhuiya/ci123chain/pkg/distribution"
	orQuery "github.com/tanhuiya/ci123chain/pkg/order"
	order "github.com/tanhuiya/ci123chain/pkg/order/rest"
	sRest "github.com/tanhuiya/ci123chain/pkg/staking/client/rest"
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
	viper.BindPFlags(rpcCmd.Flags())
}

var rpcCmd = &cobra.Command{
	Use: "rest-server",
	Short: "Start rpc server",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	order.RegisterTxRoutes(cliCtx, r)
	orQuery.RegisterTxRoutes(cliCtx, r)
	sRest.RegisterRestTxRoutes(cliCtx, r)
	sRest.RegisterTxRoutes(cliCtx, r)

	return &RestServer{
		Mux: r,
		CliCtx: cliCtx,
	}
}

const CorePrefix = "/core"

type HeightParams struct {
	Height   string   `json:"height"`
}

type QueryParams struct {
	Data    HeightParams  `json:"data"`
}

type Response struct {
	Ret 	uint32 	`json:"ret"`
	Data 	interface{}	`json:"data"`
	Message	string	`json:"message"`
}

func Handle404() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		//req.RequestURI = ""
		cli := &http.Client{}

		nodeUri := req.RequestURI
		if strings.HasPrefix(nodeUri, CorePrefix) {
			arr := strings.SplitAfter(nodeUri, CorePrefix)
			arr = arr[1:]
			newPath := strings.Join(arr, "")
			dest := viper.GetString(helper.FlagNode)
			dest = strings.ReplaceAll(dest, "tcp", "http")

			req.ParseForm()
			var data = map[string]string{}

			for k, v := range req.Form {
				//fmt.Println("key is: ", k)
				//fmt.Println("val is: ", v)
				//data.Set(k, v[0])
				key := k
				value := v[0]
				data[key] = value
			}

			newData := url.Values{}
			for k, v := range data {
				//fmt.Println(j)
				newData.Set(k, v)
			}
/*
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				//
			}
			var p QueryParams
			err = json.Unmarshal(body, &p)
			*/

			proxyurl, _ := url.Parse(dest)
			//proxy := httputil.NewSingleHostReverseProxy(proxyurl)


			remote_addr := "http://" + proxyurl.Host + newPath

			r, Err := http.NewRequest(req.Method, remote_addr, strings.NewReader(newData.Encode()))
			if Err != nil {
				panic(Err)
			}
			r.Body = ioutil.NopCloser(strings.NewReader(newData.Encode()))



			r.URL.Host = proxyurl.Host
			r.URL.Path = newPath
			//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			//req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
			//req.RequestURI = newPath
			/*
			r.URL.Host = proxyurl.Host
			r.URL.Path = newPath
			//r.RequestURI = newPath

			r.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
			//req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			// Note that ServeHttp is non blocking and uses a go routine under the hood
*/
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rep, err := cli.Do(r)
			if err != nil {
				//return nil, nil, err
			}
			resBody, err := ioutil.ReadAll(rep.Body)

			var tmResponse client.TMResponse
			err = json.Unmarshal(resBody, &tmResponse)
			ResultResponse := client.Response{
				Ret:     0,
				Data:    tmResponse,
				Message: "",
			}

			ResultByte, _ := json.Marshal(ResultResponse)

			w.Header().Set("Content-Type","application/json")
			w.Write(ResultByte)

			//w.Header().Set("Content-Type","application/json")
			//proxy.ServeHTTP(w, req)
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

