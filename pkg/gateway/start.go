package gateway

import (
	"flag"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/gateway/backend"
	"github.com/tanhuiya/ci123chain/pkg/gateway/couchdbsource"
	"github.com/tanhuiya/ci123chain/pkg/gateway/logger"
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
	var statedb, dbname, urlreg string
	var port int
	flag.StringVar(&logDir, "logdir", DefaultLogDir, "log dir")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "level for log")
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.StringVar(&statedb, "statedb", "couchdb://couchdb_service:5984", "server resource")
	flag.StringVar(&urlreg, "urlreg", "http://***:80", "reg for url connection to node")

	flag.StringVar(&dbname, "db", "ci123", "db name")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	if ok, err :=  regexp.MatchString("[*]+", urlreg); !ok {
		panic(err)
	}
	// 初始化logger
	logger.Init(logDir, "gateway", "", logLevel)

	svr := couchdbsource.NewCouchSource(dbname, statedb, urlreg)

	serverPool = NewServerPool(backend.NewBackEnd, svr, 10)

	list := strings.Split(serverList, ",")
	serverPool.ConfigServerPool(list)

	serverPool.Run()
	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(AllHandle),
	}

	// start health checking
	go healthCheck()

	go fetchSharedRoutine()

	logger.Info("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func AllHandle(w http.ResponseWriter, r *http.Request) {
	//do something
	job := NewSpecificJob(w, r, serverPool.backends)
	if job != nil {
		serverPool.JobQueue <- job
	}
	time.Sleep(100 * time.Millisecond)
}