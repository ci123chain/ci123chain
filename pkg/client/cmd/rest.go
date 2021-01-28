package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	accountRpc "github.com/ci123chain/ci123chain/pkg/account/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/cmd/rpc"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	txRpc "github.com/ci123chain/ci123chain/pkg/transfer/rest"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	dRest "github.com/ci123chain/ci123chain/pkg/distribution/client/rest"
	iRest "github.com/ci123chain/ci123chain/pkg/infrastructure/client/rest"
	mRest "github.com/ci123chain/ci123chain/pkg/mint/client/rest"
	orQuery "github.com/ci123chain/ci123chain/pkg/order"
	order "github.com/ci123chain/ci123chain/pkg/order/rest"
	sRest "github.com/ci123chain/ci123chain/pkg/staking/client/rest"
	wRest "github.com/ci123chain/ci123chain/pkg/vm/client/rest"
	"github.com/gorilla/mux"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ltypes "github.com/tendermint/tendermint/rpc/lib/types"
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
	r.HandleFunc("/healthcheck", HealthCheckHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/exportLog", ExportLogHandler(cliCtx)).Methods("GET")
	rpc.RegisterRoutes(cliCtx, r)
	accountRpc.RegisterRoutes(cliCtx, r)
	txRpc.RegisterTxRoutes(cliCtx, r)
	ibc.RegisterRoutes(cliCtx, r)
	dRest.RegisterRoutes(cliCtx, r)
	order.RegisterTxRoutes(cliCtx, r)
	orQuery.RegisterTxRoutes(cliCtx, r)
	sRest.RegisterRoutes(cliCtx, r)
	wRest.RegisterRoutes(cliCtx, r)
	mRest.RegisterRoutes(cliCtx, r)
	iRest.RegisterRoutes(cliCtx, r)

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

type ResBlock ctypes.ResultBlock

type Res struct {
	Ret 	uint32 	`json:"ret"`
	Data 	ctypes.ResultBlock	`json:"data"`
	Message	string	`json:"message"`
}

func Handle404() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		//req.RequestURI = ""
		cli := &http.Client{}
		nodeUri := req.RequestURI
		if !strings.HasPrefix(nodeUri, CorePrefix){
			http.Error(w, "404 path not found", http.StatusNotFound)
			return
		}

		arr := strings.SplitAfter(nodeUri, CorePrefix)
		arr = arr[1:]
		newPath := strings.Join(arr, "")
		dest := viper.GetString(helper.FlagNode)
		dest = strings.ReplaceAll(dest, "tcp", "http")

		req.ParseForm()
		var data = map[string]string{}

		for k, v := range req.Form {
			key := k
			value := v[0]
			data[key] = value
		}

		newData := url.Values{}
		for k, v := range data {
			newData.Set(k, v)
		}

		proxyurl, _ := url.Parse(dest)

		remote_addr := "http://" + proxyurl.Host + newPath

		r, Err := http.NewRequest(req.Method, remote_addr, strings.NewReader(newData.Encode()))
		if Err != nil {
			panic(Err)
		}
		r.Body = ioutil.NopCloser(strings.NewReader(newData.Encode()))

		r.URL.Host = proxyurl.Host
		r.URL.Path = newPath

		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rep, err := cli.Do(r)
		if err != nil  || rep.StatusCode != http.StatusOK {
			http.Error(w, rep.Status, rep.StatusCode)
			return
		}

		resBody, err := ioutil.ReadAll(rep.Body)
		var resultResponse client.Response
		if err != nil {
			resultResponse = client.Response{
				Ret:     -1,
				Data:    nil,
				Message: err.Error(),
			}
		}

		var tmResponse ltypes.RPCResponse
		err = json.Unmarshal(resBody, &tmResponse)
		if err != nil {
			resultResponse = client.Response{
				Ret:     -1,
				Data:    nil,
				Message: err.Error(),
			}
		}else {
			if tmResponse.Result == nil {
				if tmResponse.Error == nil {
					resultResponse = client.Response{
						Ret:     1,
						Data:    nil,
						Message: "response is empty",
					}
				}else {
					resultResponse = client.Response{
						Ret:     1,
						Data:    nil,
						Message: tmResponse.Error.Error(),
					}
				}
			}else {
				resStr := string(tmResponse.Result)
				var result interface{}
				switch newPath {
				case "/block":
					var a ctypes.ResultBlock
					err = cdc.UnmarshalJSON([]byte(resStr), &a)
					if err != nil {
						resultResponse = client.Response{
							Ret:     -1,
							Data:    nil,
							Message: err.Error(),
						}
					}else {
						result = a
					}
				case "/status":
					var a ctypes.ResultStatus
					err = cdc.UnmarshalJSON([]byte(resStr), &a)
					if err != nil {
						resultResponse = client.Response{
							Ret:     -1,
							Data:    nil,
							Message: err.Error(),
						}
					}else {
						result = a
					}
				case "/tx":
					var a ctypes.ResultTx
					err = cdc.UnmarshalJSON([]byte(resStr), &a)
					if err != nil {
						resultResponse = client.Response{
							Ret:     -1,
							Data:    nil,
							Message: err.Error(),
						}
					}else {
						result = a
					}
				case "/health":
					var a ctypes.ResultHealth
					err = cdc.UnmarshalJSON([]byte(resStr), &a)
					if err != nil {
						resultResponse = client.Response{
							Ret:     -1,
							Data:    nil,
							Message: err.Error(),
						}
					}else {
						result = a
					}
				}

				resultResponse = client.Response{
					Ret:     1,
					Data:    result,
					Message: "",
				}
			}
		}

		resultByte, err := json.Marshal(resultResponse)
		if err != nil {
			http.Error(w, "http response marshal error " + err.Error(), rep.StatusCode)
			return
		}

		w.Header().Set("Content-Type","application/json")
		_, _ = w.Write(resultByte)

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

func HealthCheckHandler(ctx context.Context) http.HandlerFunc  {
	return func(w http.ResponseWriter, req *http.Request) {
		cli := &http.Client{}

		dest := viper.GetString(helper.FlagNode)
		dest = strings.ReplaceAll(dest, "tcp", "http")

		req.ParseForm()
		var data = map[string]string{}

		for k, v := range req.Form {
			key := k
			value := v[0]
			data[key] = value
		}

		newData := url.Values{}
		for k, v := range data {
			newData.Set(k, v)
		}
		path := "/status"
		proxyurl, _ := url.Parse(dest)

		remote_addr := "http://" + proxyurl.Host + path

		r, Err := http.NewRequest(req.Method, remote_addr, strings.NewReader(newData.Encode()))
		if Err != nil {
			panic(Err)
		}
		r.Body = ioutil.NopCloser(strings.NewReader(newData.Encode()))

		r.URL.Host = proxyurl.Host
		r.URL.Path = path

		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rep, err := cli.Do(r)
		if err != nil  || rep.StatusCode != http.StatusOK {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(500)
			var data interface{}
			if err != nil {
				data = err
			} else {
				data = rep.Status
			}
			resultResponse := client.HealthcheckResponse{
				State:   500,
				Data:    data,
			}

			resultByte, _ := json.Marshal(resultResponse)
			w.Write(resultByte)
			return
		}

		resBody, err := ioutil.ReadAll(rep.Body)

		var tmResponse client.TMResponse
		err = json.Unmarshal(resBody, &tmResponse)
		resultResponse := client.HealthcheckResponse{
			State:   200,
			Data:    tmResponse,
		}

		resultByte, _ := json.Marshal(resultResponse)

		w.Header().Set("Content-Type","application/json")
		w.Write(resultByte)
	}
}

func ExportLogHandler(ctx context.Context) http.HandlerFunc  {
	return func(w http.ResponseWriter, req *http.Request) {
		logPath := req.FormValue("path")
		log, err := ioutil.ReadFile(logPath)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		resp := base64.StdEncoding.EncodeToString(log)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(resp))
	}
}

