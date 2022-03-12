package gateway

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/gateway/backend"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/redissource"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/pretty66/gosdk/cienv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultLogDir = "$HOME/.gateway"
	flagCiStateDBType  = "statedb_type"
	flagCiStateDBHost  = "statedb_host"
	flagCiStateDBPort  = "statedb_port"
	flagCiStateDBTls   = "statedb_tls"
	flagRPCPort 	   = "port_tm"
	flagShard          = "port_shard"
	flagETHRPCPort     = "port_eth"
	AppID			   = "hedlzgp1u48kjf50xtcvwdklminbqe9a"
	flagMaxConnection  = "max_connection"
	flagMode		   = "mode"
	flagServerList	   = "server_list"
)

var serverPool *ServerPool

var pubsubRoom *types.PubSubRoom

func Start() {
	var logLevel, serverList string
	var statedb string
	var port int
	flag.String("logdir", DefaultLogDir, "log dir")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "level for log")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.String(flagServerList, "", "Load balanced backends, use commas to separate")

	flag.String(flagMode, "lite", "gateway run mode")
	flag.String(flagRPCPort, "443", "tendermint port for websocket")
	flag.String(flagShard, "443", "shard port for websocket")
	flag.String(flagETHRPCPort, "443", "eth port for websocket")
	flag.String(flagCiStateDBType, "redis", "database types")
	flag.String(flagCiStateDBHost, "", "db host")
	flag.Uint64(flagCiStateDBPort, 7443, "db port")
	flag.Bool(flagCiStateDBTls, true, "use tls")
	flag.Int(flagMaxConnection, 20, "proxy max connection")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.SetEnvPrefix("CI")
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	runMode := viper.GetString(flagMode)
	if runMode == "" {
		runMode = "lite"
	}

	dbType := viper.GetString(flagCiStateDBType)
	if dbType == "" {
		dbType = "redis"
	}

	svr := types.ServerSource(nil)
	switch runMode {
	case "lite":
		break
	default:
		dbHost := viper.GetString(flagCiStateDBHost)
		if dbHost == "" {
			var err error
			dbHost, err = util.GetDomain()
			if err != nil {
				logger.Error("get remote db host failed", "err", err.Error())
				return
			}
		}
		if dbHost == "" {
			panic(errors.New(fmt.Sprintf("%s can not be empty", "ci-statedb-host")))
		}

		dbTls := viper.GetBool(flagCiStateDBTls)
		dbPort := viper.GetUint64(flagCiStateDBPort)
		p := strconv.FormatUint(dbPort, 10)
		switch dbType {
		case "redis":
			statedb = "redisdb://" + dbHost + ":" + p
			if dbTls {
				statedb += "#tls"
			}
		default:
			panic(errors.New(fmt.Sprintf("types of db which is not reids not implement yet")))
		}
		svr = redissource.NewRedisSource(statedb)
	}

	//logDir = viper.GetString("logdir")
	port = viper.GetInt("port")
	tmport := viper.GetString(flagRPCPort)
	shardport := viper.GetString(flagShard)
	ethport := viper.GetString(flagETHRPCPort)

	maxConnections := viper.GetInt(flagMaxConnection)
	if maxConnections < 1 {
		panic(errors.New(fmt.Sprintf("invalid max_connections: %v", maxConnections)))
	}

	// 初始化logger
	//logger.Init()
	//dynamic.Init()
	//init PubSubRoom
	pubsubRoom = &types.PubSubRoom{}
	types.SetDefaultPort(tmport,shardport, ethport)
	pubsubRoom.GetPubSubRoom()
	pubsubRoom.MaxConnections = maxConnections

	serverPool = NewServerPool(backend.NewBackEnd, svr, 10)
	serverList = viper.GetString(flagServerList)
	list := strings.Split(serverList, ",")
	serverPool.ConfigServerPool(list)

	serverPool.Run()
	// create http server

	timeoutHandler := http.TimeoutHandler(http.HandlerFunc(AllHandle), time.Second*60, "server timeout")
	http.HandleFunc("/healthcheck", healthCheckHandlerFn)
	http.Handle("/", timeoutHandler)
	http.HandleFunc("/pubsub", PubSubHandle)
	http.HandleFunc("/eth/pubsub", EthPubSubHandle)
	http.HandleFunc("/getAppkeyChannel", appKeyChannelHandlerFn)
	http.HandleFunc("/getGatewayDomain", gatewayDomainHandlerFn)

	// start health checking
	go healthCheck()

	if runMode != "lite" {
		go fetchSharedRoutine()
	}

	//check pubsub backends.
	go checkBackend()

	//check eth pubsub
	go checkEthBackend()

	logger.Info("Load Balancer started at :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Info("info: %s", err)
	}
}

func AllHandle(w http.ResponseWriter, r *http.Request) {
	//do something
	w.Header().Set("Content-Type", "application/json")
	job := NewSpecificJob(r, serverPool.backends)
	if job != nil {
		serverPool.JobQueue <- job
	} else {
		err := errors.New("arguments error")
		res, _ := json.Marshal(types.ErrorResponse{
			Ret:     0,
			Message: err.Error(),
		})
		_, _ = w.Write(res)
		return
	}
	select {
	case resp := <-*job.ResponseChan:
		_, _ = w.Write(resp)
	}
}

func healthCheckHandlerFn(w http.ResponseWriter, _ *http.Request) {
	deadList := serverPool.getDeadList()
	if len(deadList) != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		resultResponse := client.HealthcheckResponse{
			State: 500,
			Data:  deadList,
		}

		resultByte, _ := json.Marshal(resultResponse)
		_, _ = w.Write(resultByte)
		return
	}
	resultResponse := client.HealthcheckResponse{
		State: 200,
		Data:  "all backends health or no backends now",
	}
	resultByte, _ := json.Marshal(resultResponse)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resultByte)
}

func appKeyChannelHandlerFn(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	res := map[string]interface{}{
		"appid":   AppID,
		"appkey":  params.Get("appkey"),
		"channel": params.Get("channel"),
		"site": map[string]interface{}{
			"cluster_id": cienv.GetEnv("IDG_SITEUID"),
			"site_id":    cienv.GetEnv("IDG_CLUSTERUID"),
		},
	}
	resultByte, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resultByte)
}

func gatewayDomainHandlerFn(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"domain":   os.Getenv("CI_GATEWAY_DOMAIN"),
	}
	resultByte, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resultByte)
}
