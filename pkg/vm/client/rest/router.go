package rest

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/websockets"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

const (
	flagWebSocket = "wsport"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
	registerApiRoutes(cliCtx, r)
}

func registerTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/vm/contract/upload", rest.MiddleHandler(cliCtx, uploadContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/init", rest.MiddleHandler(cliCtx, instantiateContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/execute", rest.MiddleHandler(cliCtx, executeContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/migrate", rest.MiddleHandler(cliCtx, migrateContractHandler, types.DefaultCodespace)).Methods("POST")

}

func registerQueryRoutes(cliCtx context.Context, r *mux.Router) {
	//r.HandleFunc("/wasm/codeSearch/list", listCodesHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/meta", queryCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/list", listContractsByCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/info", queryContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/query", queryContractStateAllHandlerFn(cliCtx)).Methods("POST")
}

func registerApiRoutes(cliCtx context.Context, r *mux.Router) {
	server := rpc.NewServer()

	apis := GetAPIs(cliCtx)

	// Register all the APIs exposed by the namespace services
	// TODO: handle allowlist and private APIs
	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}

	// Web3 RPC API route
	r.HandleFunc("/", Handler(server)).Methods("POST", "OPTIONS")

	websocketAddr := viper.GetString(flagWebSocket)
	ws := websockets.NewServer(cliCtx, websocketAddr)
	ws.Start()
}


func Handler(s *rpc.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		w.Header().Set("content-types", "application/json;charset=UTF-8")             //返回数据格式是json
		s.ServeHTTP(w, r)
	}
}

//types Handler struct {
//	sr *http.Server
//}

//func NewHandler(sr *http.Server) Handler {
//	return Handler{sr:sr}
//}
//
//func (h Handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	println("receive not found handle request")
//	by, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		println("err:", err.Error())
//	}else {
//		println(string(by))
//	}
//}