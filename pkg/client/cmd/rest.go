package cmd

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	grpctypes "github.com/ci123chain/ci123chain/pkg/abci/types/grpc"
	accountRpc "github.com/ci123chain/ci123chain/pkg/account/rest"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/cmd/rpc"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	gravity "github.com/ci123chain/ci123chain/pkg/gravity/types"
	txRpc "github.com/ci123chain/ci123chain/pkg/transfer/rest"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	rpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	dRest "github.com/ci123chain/ci123chain/pkg/distribution/client/rest"
	ibctransferRest "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/client/rest"
	ibccore "github.com/ci123chain/ci123chain/pkg/ibc/core/client/rest"

	gRest "github.com/ci123chain/ci123chain/pkg/gravity/client/rest"
	iRest "github.com/ci123chain/ci123chain/pkg/infrastructure/client/rest"
	mRest "github.com/ci123chain/ci123chain/pkg/mint/client/rest"
	orQuery "github.com/ci123chain/ci123chain/pkg/order"
	order "github.com/ci123chain/ci123chain/pkg/order/rest"
	sRest "github.com/ci123chain/ci123chain/pkg/staking/client/rest"
	wRest "github.com/ci123chain/ci123chain/pkg/vm/client/rest"
	"github.com/gorilla/mux"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ltypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"gitlab.oneitfarm.com/bifrost/sesdk"
	"gitlab.oneitfarm.com/bifrost/sesdk/discovery"
)

const (
	FlagListenAddr         = "laddr"
	FlagMaxOpenConnections = "max-open"
	FlagRPCReadTimeout     = "read-timeout"
	FlagRPCWriteTimeout    = "write-timeout"
	FlagWebsocket		   = "wsport"
	GenesisFile			   = "genesis.json"
	PrivValidatorKey	   = "priv_validator_key.json"
	flagETHChainID         = "eth_chain_id"
	flagTokenName		   = "tokenname"
)

type ConfigFiles struct {
	GenesisFile []byte `json:"genesis_file"`
	NodeID 		string `json:"node_id"`
}

func init() {
	rootCmd.AddCommand(rpcCmd)
	rpcCmd.Flags().String(FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	rpcCmd.Flags().Uint(FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	rpcCmd.Flags().Uint(FlagRPCReadTimeout, 10, "The RPC read timeout")
	rpcCmd.Flags().Uint(FlagRPCWriteTimeout, 10, "The RPC write timeout")
	rpcCmd.Flags().String(FlagWebsocket, "8546", "websocket port to listen to")
	rpcCmd.Flags().Int64(flagETHChainID, 1, "eth_chain_id")
	rpcCmd.Flags().String(flagTokenName, "stake", "Chain token name")

	_ = viper.BindPFlags(rpcCmd.Flags())
}

var rpcCmd = &cobra.Command{
	Use: "rest-server",
	Short: "Start rpc server",
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetInt64(flagETHChainID)
		util.Setup(id)
		denom := viper.GetString(flagTokenName)
		types.SetCoinDenom(denom)
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
	GRPCGatewayRouter *runtime.ServeMux
	CliCtx  context.Context
	listener net.Listener
}

func NewRestServer() *RestServer {
	r := mux.NewRouter()
	cliCtx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return nil
	}
	go SetupRegisterCenter()

	r.NotFoundHandler = Handle404()
	r.HandleFunc("/healthcheck", HealthCheckHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/exportLog", ExportLogHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/exportConfig", ExportConfigHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/exportEnv", ExportEnv(cliCtx)).Methods("POST")
	r.HandleFunc("/info", registerCenterHandler(cliCtx)).Methods("GET")
	rpc.RegisterRoutes(cliCtx, r)
	accountRpc.RegisterRoutes(cliCtx, r)
	txRpc.RegisterTxRoutes(cliCtx, r)
	// todo ibc
	dRest.RegisterRoutes(cliCtx, r)
	order.RegisterTxRoutes(cliCtx, r)
	orQuery.RegisterTxRoutes(cliCtx, r)
	sRest.RegisterRoutes(cliCtx, r)
	wRest.RegisterRoutes(cliCtx, r)
	mRest.RegisterRoutes(cliCtx, r)
	iRest.RegisterRoutes(cliCtx, r)
	gRest.RegisterRoutes(cliCtx, r, gravity.StoreKey)


	ibctransferRest.RegisterRoutes(cliCtx, r)
	ibccore.RegisterRoutes(cliCtx, r)

	// The default JSON marshaller used by the gRPC-Gateway is unable to marshal non-nullable non-scalar fields.
	// Using the gogo/gateway package with the gRPC-Gateway WithMarshaler option fixes the scalar field marshalling issue.
	marshalerOption := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
		AnyResolver:  cliCtx.InterfaceRegistry,
	}

	grpcRoute := runtime.NewServeMux(
		// Custom marshaler option is required for gogo proto
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshalerOption),

		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),

		// Custom header matcher for mapping request headers to
		// GRPC metadata
		runtime.WithIncomingHeaderMatcher(CustomGRPCHeaderMatcher),
	)
	module.ModuleBasics.RegisterGRPCGatewayRoutes(cliCtx, grpcRoute)

	return &RestServer{
		Mux: r,
		GRPCGatewayRouter: grpcRoute,
		CliCtx: cliCtx,
	}
}

// CustomGRPCHeaderMatcher for mapping request headers to
// GRPC metadata.
// HTTP headers that start with 'Grpc-Metadata-' are automatically mapped to
// gRPC metadata after removing prefix 'Grpc-Metadata-'. We can use this
// CustomGRPCHeaderMatcher if headers don't start with `Grpc-Metadata-`
func CustomGRPCHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case grpctypes.GRPCBlockHeightHeader:
		return grpctypes.GRPCBlockHeightHeader, true
	default:
		return runtime.DefaultHeaderMatcher(key)
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

		_ = req.ParseForm()
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
						Message: tmResponse.Error.Message,
					}
				}
			}else {
				resultResponse = client.Response{
					Ret:     1,
					Data:    tmResponse.Result,
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
		err = rs.listener.Close()
		return
	})

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	//rs.registerGRPCGatewayRoutes()

	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	return rpcserver.Serve(rs.listener, rs.Mux, logger,cfg)
}

func (rs *RestServer) registerGRPCGatewayRoutes() {
	rs.Mux.PathPrefix("/").Handler(rs.GRPCGatewayRouter)
}

func HealthCheckHandler(ctx context.Context) http.HandlerFunc  {
	return func(w http.ResponseWriter, req *http.Request) {
		cli := &http.Client{}

		dest := viper.GetString(helper.FlagNode)
		dest = strings.ReplaceAll(dest, "tcp", "http")

		_ = req.ParseForm()
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
		viper.SetEnvPrefix("CI")
		_ = viper.BindEnv("HOME")
		root := viper.GetString("HOME")
		logPath := req.URL.Query().Get("path")
		logger, err := ioutil.ReadFile(filepath.Join(root, "logs", logPath))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", logPath))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(logger)
	}
}

func ExportConfigHandler(ctx context.Context) http.HandlerFunc  {
	return func(w http.ResponseWriter, req *http.Request) {
		viper.SetEnvPrefix("CI")
		_ = viper.BindEnv("HOME")
		root := viper.GetString("HOME")
		gen, err := ioutil.ReadFile(filepath.Join(root, "config", GenesisFile))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		pv, err := ioutil.ReadFile(filepath.Join(root, "config", PrivValidatorKey))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		var key *privval.FilePVKey
		err = cdc.UnmarshalJSON(pv, &key)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		nodeID := strings.ToLower(key.Address.String())
		configFile := &ConfigFiles{
			GenesisFile: gen,
			NodeID:      nodeID,
		}
		configBytes, err := json.Marshal(configFile)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(configBytes)
	}
}

func ExportEnv(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var ks []string
		keys := req.FormValue("keys")
		err := json.Unmarshal([]byte(keys), &ks)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		var res = make(map[string]interface{}, 0)
		for _, v := range ks {
			value := os.Getenv(v)
			res[v] = value
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(bytes)
	}
}

func SetupRegisterCenter() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "register-center")
	appID := os.Getenv("CI_VALIDATOR_KEY")
	if appID == "" {
		logger.Error("CI_VALIDATOR_KEY can not be empty")
		os.Exit(1)
	}
	hn, _ := os.Hostname()
	// 注册中心自身，初始化配置
	conf := &discovery.Config{
		// discovery地址
		Nodes:    []string{"192.168.2.80:8181", "192.168.2.80:8182", "192.168.2.80:8183"},
		Region:   "sh",
		Zone:     "sh001",
		Env:      "dev",
		Host:     hn,               // hostname
		RenewGap: time.Second * 30, // 心跳时间
	}
	// 自身实例信息
	ins := &sesdk.Instance{
		Region:   "sh",
		Zone:     "sh001",
		Env:      "dev",
		AppID:    appID, // 自身唯一识别号
		Hostname: hn,
		Addrs: []string{ // 可上报任意服务监听地址，供发现方连接
			"http://127.0.0.1:8545",
			//"https://127.0.0.1:443",
			//"tcp://192.168.2.88:3030",
		},
		// 上报任意自身属性信息
		Metadata: map[string]string{
			"weight":       "10", // 负载均衡权重
			"runtime":      "production",
			"service_name": "tttttttttt",
		},
	}
	// 实例化discovery对象
	dis, err := discovery.New(conf)
	if err != nil {
		panic(err)
	}
	// 注册自身
	_, err = dis.Register(ctx.Background(), ins)
	if err != nil {
		panic(err)
	}
	//// 启动服务主要逻辑
	//go func() {
	//	http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
	//		_, _ = writer.Write([]byte(`{"state":1,"msg":"OK"}`))
	//	})
	//
	//	err := http.ListenAndServe(":38081", nil)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()
	// 监听系统信号，服务下线
	dis.ExitSignal(func(s os.Signal) {
		logger.Info("got exit signal, exit now", "signal", s.String())
	})
}

func registerCenterHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//var ks []string
		//keys := req.FormValue("keys")
		//err := json.Unmarshal([]byte(keys), &ks)
		//if err != nil {
		//	_, _ = w.Write([]byte(err.Error()))
		//	return
		//}
		//var res = make(map[string]interface{}, 0)
		//for _, v := range ks {
		//	value := os.Getenv(v)
		//	res[v] = value
		//}
		res := map[string]interface{}{
			"state": 200,
			"appID": os.Getenv("CI_VALIDATOR_KEY"),
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(bytes)
	}
}