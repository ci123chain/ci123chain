package gateway

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/gateway/backend"
	"github.com/ci123chain/ci123chain/pkg/gateway/couchdbsource"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const DefaultLogDir = "$HOME/.gateway"

var serverPool *ServerPool

var pubsubRoom *types.PubSubRoom

func Start() {
	var logLevel, serverList string
	var statedb, urlreg string
	var port int
	flag.String("logdir", DefaultLogDir, "log dir")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "level for log")

	flag.StringVar(&serverList, "backends", "http://localhost:1317", "Load balanced backends, use commas to separate")
	flag.String("statedb", "couchdb://couchdb_service:5984/ci123", "server resource")
	flag.StringVar(&urlreg, "urlreg", "http://***:80", "reg for url connection to node")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.String("rpcport", "80", "rpc address for websocket")
	//flag.Parse()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.SetEnvPrefix("CI")
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	//viper.BindEnv("statedb")
	//viper.BindEnv("logdir")
	statedb = viper.GetString("statedb")
	//logDir = viper.GetString("logdir")
	port = viper.GetInt("port")
	urlreg = viper.GetString("urlreg")
	rpcAddress := viper.GetString("rpcport")
	if rpcAddress == "" {
		rpcAddress = types.DefaultRPCPort
	}

	if ok, err := regexp.MatchString("[*]+", urlreg); !ok {
		panic(err)
	}
	// 初始化logger
	logger.Init()
	//dynamic.Init()
	//init PubSubRoom
	pubsubRoom = &types.PubSubRoom{}
	types.SetDefaultPort(rpcAddress)
	pubsubRoom.GetPubSubRoom()

	svr := couchdbsource.NewCouchSource(statedb, urlreg)

	serverPool = NewServerPool(backend.NewBackEnd, svr, 10)

	list := strings.Split(serverList, ",")
	serverPool.ConfigServerPool(list)

	serverPool.Run()
	// create http server

	timeoutHandler := http.TimeoutHandler(http.HandlerFunc(AllHandle), time.Second*60, "server timeout")
	http.HandleFunc("/healthcheck", healthCheckHandlerFn)
	http.Handle("/", timeoutHandler)
	http.HandleFunc("/pubsub", PubSubHandle)

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
