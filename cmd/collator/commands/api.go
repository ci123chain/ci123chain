package commands

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	helpers "github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"runtime"
	"sync"
)


// Service represents a relayer listen service
// TODO: sync services to disk so that they can survive restart
type Service struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Src    string `json:"src"`
	//SrcKey string `json:"src-key"`
	Dst    string `json:"dst"`
	//DstKey string `json:"dst-key"`

	doneFunc func()
}

// NewService returns a new instance of Service
func NewService(name, path string, src, dst *collactor.Chain, doneFunc func()) *Service {
	return &Service{name, path, src.ChainID, dst.ChainID,  doneFunc}
}

// ServicesManager represents the manager of the various services the relayer is running
type ServicesManager struct {
	Services map[string]*Service

	sync.Mutex
}

// NewServicesManager returns a new instance of a services manager
func NewServicesManager() *ServicesManager {
	return &ServicesManager{Services: make(map[string]*Service)}
}

func getAPICmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use: "api",
		Short: "Start the relayer API",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := mux.NewRouter()

			sm := NewServicesManager()
			// VERSION
			// version get
			r.HandleFunc("/version", VersionHandler).Methods("GET")
			// CONFIG
			// config get
			r.HandleFunc("/config", ConfigHandler).Methods("GET")
			// CHAINS
			// chains get
			r.HandleFunc("/chains", GetChainsHandler).Methods("GET")
			// chain get
			r.HandleFunc("/chains/{name}", GetChainHandler).Methods("GET")
			// chain status get
			r.HandleFunc("/chains/{name}/status", GetChainStatusHandler).Methods("GET")
			// chain add
			// NOTE: when updating the config, we need to update the global config object ()
			// as well as the file on disk updating the file on disk may cause contention, just
			// retry the `file.Open` until the file can be opened, rewrite it with the new config
			// and close it
			r.HandleFunc("/chains/{name}", PostChainHandler).Methods("POST")
			// chain update
			//r.HandleFunc("/chains/{name}", PutChainHandler).Methods("PUT")
			// chain delete
			r.HandleFunc("/chains/{name}", DeleteChainHandler).Methods("DELETE")

			// PATHS
			// paths get
			r.HandleFunc("/paths", GetPathsHandler).Methods("GET")
			// path get
			r.HandleFunc("/paths/{name}", GetPathHandler).Methods("GET")
			// path status get
			r.HandleFunc("/paths/{name}/status", GetPathStatusHandler).Methods("GET")
			// path add
			r.HandleFunc("/paths/{name}", PostPathHandler).Methods("POST")
			// path delete
			r.HandleFunc("/paths/{name}", DeletePathHandler).Methods("DELETE")

			// LIGHT
			// light header, if no ?height={height} is passed, latest
			r.HandleFunc("/light/{chain-id}/header", GetLightHeader).Methods("GET")
			// light height
			r.HandleFunc("/light/{chain-id}/height", GetLightHeight).Methods("GET")
			// light create
			r.HandleFunc("/light/{chain-id}", PostLight).Methods("POST")
			// light update
			//r.HandleFunc("/light/{chain-id}", PutLight).Methods("PUT")
			// light delete
			r.HandleFunc("/light/{chain-id}", DeleteLight).Methods("DELETE")

			r.HandleFunc("/link/{path}", PostLinkChain).Methods("POST")

			// TODO: this particular function needs some work, we need to listen on chains in configuration and
			// route all the events (both block and tx) though and event bus to allow for multiple subscribers
			// on update of config we need to handle that case
			// Data for this should be stored in the ServicesManager struct
			r.HandleFunc("/listen", GetRelayerListenHandler(sm)).Methods("GET")
			r.HandleFunc("/listen/{path}/{strategy}/{name}", PostRelayerListenHandler(sm)).Methods("POST")
			r.HandleFunc("/listen/{path}/{strategy}/{name}", DeleteRelayerListenHandler(sm)).Methods("DELETE")


			fmt.Println("listening on", config.Global.APIListenPort)

			if err := http.ListenAndServe(config.Global.APIListenPort, r); err != nil {
				return err
			}

			return nil
		},
	}
	return apiCmd
}


// VersionHandler returns the version info in json format
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	version := versionInfo{
		Version:   Version,
		Commit:    Commit,
		Chain: 	   SDKCommit,
		Go:        fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
	helpers.SuccessJSONResponse(http.StatusOK, version, w)
}

// ConfigHandler handles the route
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	helpers.SuccessJSONResponse(http.StatusOK, config, w)
}

// API Handlers

// GetChainsHandler returns the configured chains in json format
func GetChainsHandler(w http.ResponseWriter, r *http.Request) {
	helpers.SuccessJSONResponse(http.StatusOK, config.Chains, w)
}

// GetChainHandler returns the configured chains in json format
func GetChainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, chain, w)
}

// GetChainStatusHandler returns the configured chains in json format
func GetChainStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, chainStatusResponse{}.Populate(chain), w)
}


// PostRelayerListenHandler returns a handler for a listener that can listen on many IBC paths
func PostRelayerListenHandler(sm *ServicesManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, found := sm.Services[vars["name"]]
		if found {
			helpers.WriteErrorResponse(http.StatusBadRequest, errors.New("relay service already exist"), w)
			return
		}
		// TODO: make this handler accept a json post argument
		pth, err := config.Paths.Get(vars["path"])
		if err != nil {
			helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
			return
		}
		c, src, dst, err := config.ChainsFromPath(vars["path"])
		if err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
		pth.Strategy = &collactor.StrategyCfg{Type: vars["strategy"]}
		strategyType, err := pth.GetStrategy()
		if err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
		done, err := collactor.RunStrategy(c[src], c[dst], strategyType)
		if err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
		sm.Lock()
		sm.Services[vars["name"]] = NewService(vars["name"], vars["path"], c[src], c[dst], done)
		sm.Unlock()

		helpers.SuccessJSONResponse(http.StatusOK, "relayer start successful", w)

	}
}



// DeleteRelayerListenHandler returns a handler for a listener that can listen on many IBC paths
func DeleteRelayerListenHandler(sm *ServicesManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// TODO: check name to ensure that no other services exist
		service, found := sm.Services[vars["name"]]
		if !found {
			helpers.WriteErrorResponse(http.StatusBadRequest, errors.New("relay service not exist"), w)
			return
		}
		if service.Path != vars["path"] {
			helpers.WriteErrorResponse(http.StatusBadRequest, errors.New("parameter path not match"), w)
			return
		}
		service.doneFunc()
		sm.Lock()
		delete(sm.Services,vars["name"])
		sm.Unlock()

		helpers.SuccessJSONResponse(http.StatusOK, "relayer delete successful", w)
	}
}

func GetRelayerListenHandler(sm *ServicesManager)  func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SuccessJSONResponse(http.StatusOK, sm.Services, w)
	}
}