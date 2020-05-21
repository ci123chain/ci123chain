package gateway

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/dynamic"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/gateway/backend"
	"github.com/ci123chain/ci123chain/pkg/gateway/couchdbsource"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const DefaultLogDir  = "$HOME/.gateway"
var serverPool *ServerPool

func Start() {
	var logDir, logLevel, serverList string
	var statedb, urlreg string
	var port int
	flag.String("logdir", DefaultLogDir, "log dir")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "level for log")

	flag.StringVar(&serverList, "backends", "http://localhost:1317", "Load balanced backends, use commas to separate")
	flag.String("statedb", "couchdb://couchdb_service:5984/ci123", "server resource")
	flag.StringVar(&urlreg, "urlreg", "http://***:80", "reg for url connection to node")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	//flag.Parse()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.SetEnvPrefix("CI")
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	//viper.BindEnv("statedb")
	//viper.BindEnv("logdir")
	statedb = viper.GetString("statedb")
	logDir = viper.GetString("logdir")

	if ok, err :=  regexp.MatchString("[*]+", urlreg); !ok {
		panic(err)
	}
	// 初始化logger
	logger.Init(logDir, "gateway", "", logLevel)

	svr := couchdbsource.NewCouchSource(statedb, urlreg)

	serverPool = NewServerPool(backend.NewBackEnd, svr, 10)

	list := strings.Split(serverList, ",")
	serverPool.ConfigServerPool(list)

	serverPool.Run()
	// create http server

	timeoutHandler := http.TimeoutHandler(http.HandlerFunc(AllHandle), time.Second*60, "server timeout")
	http.HandleFunc("/createChannel", createChannelHandle)
	http.HandleFunc("/testCall", testCallHandle)
	http.Handle("/", timeoutHandler)

	// start health checking
	go healthCheck()

	go fetchSharedRoutine()

	logger.Info("Load Balancer started at :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func createChannelHandle(w http.ResponseWriter, r *http.Request) {
	err := dynamic.CreateChannel(dynamic.APP_ID, dynamic.APP_KEY, dynamic.CHANNEL,"", r.Header)
	if err != nil {
		log.Fatal(err)
	}
}

func testCallHandle(w http.ResponseWriter, r *http.Request) {
	err := dynamic.TestCall(dynamic.APP_ID, dynamic.APP_KEY, dynamic.CHANNEL,"", r.Header)
	if err != nil {
		log.Fatal(err)
	}
}

func AllHandle(w http.ResponseWriter, r *http.Request) {
	//do something
	w.Header().Set("Content-Type", "application/json")
	job := NewSpecificJob(r, serverPool.backends)
	if job != nil {
		serverPool.JobQueue <- job
	}else {
		err := errors.New("arguments error")
		res, _ := json.Marshal(types.ErrorResponse{
			Err:  err.Error(),
		})
		_, _ = w.Write(res)
	}
	select {
	 case resp := <- *job.ResponseChan:
		_, _ = w.Write(resp)
	}
}