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
	"regexp"
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
	flagRPCPort 	   = "rpcport"
	AppID			   = "hedlzgp1u48kjf50xtcvwdklminbqe9a"
)

var serverPool *ServerPool

var pubsubRoom *types.PubSubRoom

func Start() {
	var logLevel, serverList string
	var statedb, urlreg string
	var port int
	flag.String("logdir", DefaultLogDir, "log dir")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "level for log")

	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.StringVar(&urlreg, "urlreg", "http://***", "reg for url connection to node")
	flag.IntVar(&port, "port", 3030, "Port to serve")

	flag.String(flagRPCPort, "80", "rpc address for websocket")
	flag.String(flagCiStateDBType, "redis", "database types")
	flag.String(flagCiStateDBHost, "", "db host")
	flag.Uint64(flagCiStateDBPort, 7443, "db port")
	flag.Bool(flagCiStateDBTls, true, "use tls")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.SetEnvPrefix("CI")
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	//viper.BindEnv("statedb")
	//viper.BindEnv("logdir")

	dbType := viper.GetString(flagCiStateDBType)
	if dbType == "" {
		dbType = "redis"
	}
	//dbHost := viper.GetString(flagCiStateDBHost)
	dbHost := util.Discovery(util.CallBack)
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
	//logDir = viper.GetString("logdir")
	port = viper.GetInt("port")
	urlreg = viper.GetString("urlreg")
	rpcAddress := viper.GetString(flagRPCPort)
	if rpcAddress == "" {
		rpcAddress = types.DefaultRPCPort
	}

	if ok, err := regexp.MatchString("[*]+", urlreg); !ok {
		panic(err)
	}
	// 初始化logger
	//logger.Init()
	//dynamic.Init()
	//init PubSubRoom
	pubsubRoom = &types.PubSubRoom{}
	types.SetDefaultPort(rpcAddress)
	pubsubRoom.GetPubSubRoom()

	svr := redissource.NewRedisSource(statedb, urlreg)

	serverPool = NewServerPool(backend.NewBackEnd, svr, 10)

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

	go fetchSharedRoutine()
	//check pubsub backends.
	go checkBackend()

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
