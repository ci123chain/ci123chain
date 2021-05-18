package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	helpers "github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func chainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chains",
		Aliases: []string{"ch"},
		Short:   "manage chain configurations",
	}

	cmd.AddCommand(
		//chainsListCmd(),
		//chainsDeleteCmd(),
		chainsAddCmd(),
		//chainsEditCmd(),
		//chainsShowCmd(),
		//chainsAddrCmd(),
		//chainsAddDirCmd(),
	)

	return cmd
}



func chainsAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new chain to the configuration file by passing a file (-f) or url (-u), or user input",
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s chains add
$ %s ch a
$ %s chains add --file chains/ibc0.json
$ %s chains add --url http://xxxx.com/ibc0.json
`, appName, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			var out *Config

			file, url, err := getAddInputs(cmd)
			if err != nil {
				return err
			}

			switch {
			case file != "":
				if out, err = fileInputAdd(file); err != nil {
					return err
				}
			case url != "":
				if out, err = urlInputAdd(url); err != nil {
					return err
				}
			default:
				return errors.New("unsupport input method")
			}

			if err = validateConfig(out); err != nil {
				return err
			}

			return overWriteConfig(out)
		},
	}

	return chainsAddFlags(cmd)
}
func fileInputAdd(file string) (cfg *Config, err error) {
	// If the user passes in a file, attempt to read the chain config from that file
	c := &collactor.Chain{}
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(byt, c); err != nil {
		return nil, err
	}

	if err = config.AddChain(c); err != nil {
		return nil, err
	}

	return config, nil
}


// urlInputAdd validates a chain config URL and fetches its contents
func urlInputAdd(rawurl string) (cfg *Config, err error) {
	u, err := url.Parse(rawurl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return cfg, errors.New("invalid URL")
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var c *collactor.Chain
	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()
	err = d.Decode(c)
	if err != nil {
		return cfg, err
	}

	if err = config.AddChain(c); err != nil {
		return nil, err
	}
	return config, err
}
//
//func overWriteConfig(cfg *Config) (err error) {
//	cfgPath := path.Join(homePath, "configs", "configs.yaml")
//	if _, err = os.Stat(cfgPath); err == nil {
//		viper.SetConfigFile(cfgPath)
//		if err = viper.ReadInConfig(); err == nil {
//			// ensure validateConfig runs properly
//			err = validateConfig(configs)
//			if err != nil {
//				return err
//			}
//
//			// marshal the new configs
//			out, err := yaml.Marshal(cfg)
//			if err != nil {
//				return err
//			}
//
//			// overwrite the configs file
//			err = ioutil.WriteFile(viper.ConfigFileUsed(), out, 0600)
//			if err != nil {
//				return err
//			}
//
//			// set the global variable
//			configs = cfg
//		}
//	}
//	return err
//}



type chainStatusResponse struct {
	Light   bool `json:"light"`
	Path    bool `json:"path"`
	Key     bool `json:"key"`
	Balance bool `json:"balance"`
}

func (cs chainStatusResponse) Populate(c *collactor.Chain) chainStatusResponse {
	_, err := c.GetAddress()
	if err == nil {
		cs.Key = true
	}

	coins, err := c.QueryBalanceWithAddress(c.MustGetAddress())
	if err == nil && !coins.Empty() {
		cs.Balance = true
	}

	_, err = c.GetLatestLightHeader()
	if err == nil {
		cs.Light = true
	}

	for _, pth := range config.Paths {
		if pth.Src.ChainID == c.ChainID || pth.Dst.ChainID == c.ChainID {
			cs.Path = true
		}
	}
	return cs
}

type addChainRequest struct {
	PrivateKey 		string `json:"private-key"`
	//Key            string `json:"key"`
	RPCAddr        string `json:"rpc-addr"`
	AccountPrefix  string `json:"account-prefix"`
	//GasAdjustment  string `json:"gas-adjustment"`
	GasPrices      string `json:"gas-prices"`
	TrustingPeriod string `json:"trusting-period"`
	// required: false
	//FilePath string `json:"file"`
	// required: false
	URL string `json:"url"`
}
// PostChainHandler handles the route
func PostChainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainID := vars["name"]

	var request addChainRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	//if request.FilePath != "" && request.URL != "" {
	//	helpers.WriteErrorResponse(http.StatusBadRequest, errMultipleAddFlags, w)
	//	return
	//}

	var out *Config
	switch {
	//case request.FilePath != "":
	//	if out, err = fileInputAdd(request.FilePath); err != nil {
	//		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
	//		return
	//	}
	case request.URL != "":
		if out, err = urlInputAdd(request.URL); err != nil {
			helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
			return
		}
	default:
		if out, err = addChainByRequest(request, chainID); err != nil {
			helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
			return
		}
	}

	if err = validateConfig(out); err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	if err = overWriteConfig(out); err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusCreated, fmt.Sprintf("chain %s added successfully", chainID), w)
}


// DeleteChainHandler handles the route
func DeleteChainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := config.Chains.Get(vars["name"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	if err := overWriteConfig(config.DeleteChain(vars["name"])); err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, fmt.Sprintf("chain %s deleted", vars["name"]), w)
}


func addChainByRequest(request addChainRequest, chainID string) (cfg *Config, err error) {
	c := &collactor.Chain{}

	if c, err = c.Update("chain-id", chainID); err != nil {
		return nil, err
	}

	//if c, err = c.Update("key", request.Key); err != nil {
	//	return nil, err
	//}
	if c, err = c.Update("private-key", request.PrivateKey); err != nil {
		return nil, err
	}

	if c, err = c.Update("rpc-addr", request.RPCAddr); err != nil {
		return nil, err
	}

	if c, err = c.Update("account-prefix", request.AccountPrefix); err != nil {
		return nil, err
	}

	//if c, err = c.Update("gas-adjustment", request.GasAdjustment); err != nil {
	//	return nil, err
	//}

	if c, err = c.Update("gas-prices", request.GasPrices); err != nil {
		return nil, err
	}

	if c, err = c.Update("trusting-period", request.TrustingPeriod); err != nil {
		return nil, err
	}

	if err = config.AddChain(c); err != nil {
		return nil, err
	}

	return config, nil
}