package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	helpers "github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


func pathsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "paths",
		Aliases: []string{"pth"},
		Short:   "manage path configurations",
		Long: `
A path represents the "full path" or "link" for communication between two chains. This includes the client, 
connection, and channel ids from both the source and destination chains as well as the strategy to use when relaying`,
	}

	cmd.AddCommand(
		//pathsListCmd(),
		//pathsShowCmd(),
		pathsAddCmd(),
		//pathsGenCmd(),
		//pathsDeleteCmd(),
	)

	return cmd
}

func pathsAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [src-chain-id] [dst-chain-id] [path-name]",
		Aliases: []string{"a"},
		Short:   "add a path to the list of paths",
		Args:    cobra.ExactArgs(3),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s paths add ibc-0 ibc-1 demo-path
$ %s paths add ibc-0 ibc-1 demo-path --file paths/demo.json
$ %s pth a ibc-0 ibc-1 demo-path`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]
			_, err := config.Chains.Gets(src, dst)
			if err != nil {
				return fmt.Errorf("chains need to be configured before paths to them can be added: %w", err)
			}

			var out *Config
			file, err := cmd.Flags().GetString(flagFile)
			if err != nil {
				return err
			}

			if file != "" {
				if out, err = fileInputPathAdd(file, args[2]); err != nil {
					return err
				}
			} else {
				if out, err = userInputPathAdd(src, dst, args[2]); err != nil {
					return err
				}
			}

			return overWriteConfig(out)
		},
	}
	return fileFlag(cmd)
}



func fileInputPathAdd(file, name string) (cfg *Config, err error) {
	// If the user passes in a file, attempt to read the chain config from that file
	p := &collactor.Path{}
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(byt, &p); err != nil {
		return nil, err
	}

	if err = config.ValidatePath(p); err != nil {
		return nil, err
	}

	if err = config.Paths.Add(name, p); err != nil {
		return nil, err
	}

	return config, nil
}


func userInputPathAdd(src, dst, name string) (*Config, error) {
	var (
		value string
		err   error
		path  = &collactor.Path{
			Strategy: collactor.NewNaiveStrategy(),
			Src: &collactor.PathEnd{
				ChainID: src,
				Order:   "ORDERED",
			},
			Dst: &collactor.PathEnd{
				ChainID: dst,
				Order:   "ORDERED",
			},
		}
	)

	fmt.Printf("enter src(%s) client-id...\n", src)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Src.ClientID = value

	if err = path.Src.Vclient(); err != nil {
		return nil, err
	}

	fmt.Printf("enter src(%s) connection-id...\n", src)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Src.ConnectionID = value

	if err = path.Src.Vconn(); err != nil {
		return nil, err
	}

	fmt.Printf("enter src(%s) channel-id...\n", src)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Src.ChannelID = value

	if err = path.Src.Vchan(); err != nil {
		return nil, err
	}

	fmt.Printf("enter src(%s) port-id...\n", src)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Src.PortID = value

	if err = path.Src.Vport(); err != nil {
		return nil, err
	}

	fmt.Printf("enter src(%s) version...\n", src)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Src.Version = value

	if err = path.Src.Vversion(); err != nil {
		return nil, err
	}

	fmt.Printf("enter dst(%s) client-id...\n", dst)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Dst.ClientID = value

	if err = path.Dst.Vclient(); err != nil {
		return nil, err
	}

	fmt.Printf("enter dst(%s) connection-id...\n", dst)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Dst.ConnectionID = value

	if err = path.Dst.Vconn(); err != nil {
		return nil, err
	}

	fmt.Printf("enter dst(%s) channel-id...\n", dst)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Dst.ChannelID = value

	if err = path.Dst.Vchan(); err != nil {
		return nil, err
	}

	fmt.Printf("enter dst(%s) port-id...\n", dst)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Dst.PortID = value

	if err = path.Dst.Vport(); err != nil {
		return nil, err
	}

	fmt.Printf("enter dst(%s) version...\n", dst)
	if value, err = readStdin(); err != nil {
		return nil, err
	}

	path.Dst.Version = value

	if err = path.Dst.Vversion(); err != nil {
		return nil, err
	}

	if err = config.ValidatePath(path); err != nil {
		return nil, err
	}

	if err = config.Paths.Add(name, path); err != nil {
		return nil, err
	}

	return config, nil
}


// API Handlers

// GetPathsHandler returns the configured chains in json format
func GetPathsHandler(w http.ResponseWriter, r *http.Request) {
	helpers.SuccessJSONResponse(http.StatusOK, config.Paths, w)
}

// GetPathHandler returns the configured chains in json format
func GetPathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pth, err := config.Paths.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, pth, w)
}


// GetPathStatusHandler returns the configured chains in json format
func GetPathStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pth, err := config.Paths.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	c, src, dst, err := config.ChainsFromPath(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	ps := pth.QueryPathStatus(c[src], c[dst])
	helpers.SuccessJSONResponse(http.StatusOK, ps, w)
}


type postPathRequest struct {
	//FilePath string          `json:"file"`
	Src      collactor.PathEnd `json:"src"`
	Dst      collactor.PathEnd `json:"dst"`
}

// PostPathHandler handles the route
func PostPathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathName := vars["name"]

	var request postPathRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	var out *Config
	//if request.FilePath != "" {
	//	if out, err = fileInputPathAdd(request.FilePath, pathName); err != nil {
	//		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
	//		return
	//	}
	//} else {
		if out, err = addPathByRequest(request, pathName); err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
	//}

	if err = overWriteConfig(out); err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusCreated, fmt.Sprintf("path %s added successfully", pathName), w)
}



func addPathByRequest(req postPathRequest, pathName string) (*Config, error) {
	var (
		path = &collactor.Path{
			Strategy: collactor.NewNaiveStrategy(),
			Src:      &req.Src,
			Dst:      &req.Dst,
		}
	)

	if err := config.ValidatePath(path); err != nil {
		return nil, err
	}

	if err := config.Paths.Add(pathName, path); err != nil {
		return nil, err
	}

	return config, nil
}


// DeletePathHandler handles the route
func DeletePathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := config.Paths.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	cfg := config
	delete(cfg.Paths, vars["name"])

	if err = overWriteConfig(cfg); err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, fmt.Sprintf("path %s deleted", vars["name"]), w)
}