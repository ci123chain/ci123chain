package commands

import (
	"github.com/spf13/cobra"
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
		//chainsAddCmd(),
		//chainsEditCmd(),
		//chainsShowCmd(),
		//chainsAddrCmd(),
		//chainsAddDirCmd(),
	)

	return cmd
}


//
//func chainsAddCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "add",
//		Aliases: []string{"a"},
//		Short:   "Add a new chain to the configuration file by passing a file (-f) or url (-u), or user input",
//		Example: strings.TrimSpace(fmt.Sprintf(`
//$ %s chains add
//$ %s ch a
//$ %s chains add --file chains/ibc0.json
//$ %s chains add --url http://relayer.com/ibc0.json
//`, appName, appName, appName, appName)),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			var out *Config
//
//			file, url, err := getAddInputs(cmd)
//			if err != nil {
//				return err
//			}
//
//			switch {
//			case file != "":
//				if out, err = fileInputAdd(file); err != nil {
//					return err
//				}
//			case url != "":
//				if out, err = urlInputAdd(url); err != nil {
//					return err
//				}
//			default:
//				return errors.New("unsupport input method")
//			}
//
//			if err = validateConfig(out); err != nil {
//				return err
//			}
//
//			return overWriteConfig(out)
//		},
//	}
//
//	return chainsAddFlags(cmd)
//}
//func fileInputAdd(file string) (cfg *Config, err error) {
//	// If the user passes in a file, attempt to read the chain config from that file
//	c := &collactor.Chain{}
//	if _, err := os.Stat(file); err != nil {
//		return nil, err
//	}
//
//	byt, err := ioutil.ReadFile(file)
//	if err != nil {
//		return nil, err
//	}
//
//	if err = json.Unmarshal(byt, c); err != nil {
//		return nil, err
//	}
//
//	if err = config.AddChain(c); err != nil {
//		return nil, err
//	}
//
//	return config, nil
//}
//
//
//
//// urlInputAdd validates a chain config URL and fetches its contents
//func urlInputAdd(rawurl string) (cfg *Config, err error) {
//	u, err := url.Parse(rawurl)
//	if err != nil || u.Scheme == "" || u.Host == "" {
//		return cfg, errors.New("invalid URL")
//	}
//
//	resp, err := http.Get(u.String())
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//
//	var c *collactor.Chain
//	d := json.NewDecoder(resp.Body)
//	d.DisallowUnknownFields()
//	err = d.Decode(c)
//	if err != nil {
//		return cfg, err
//	}
//
//	if err = config.AddChain(c); err != nil {
//		return nil, err
//	}
//	return config, err
//}
//
//
//func overWriteConfig(cfg *Config) (err error) {
//	cfgPath := path.Join(homePath, "config", "config.yaml")
//	if _, err = os.Stat(cfgPath); err == nil {
//		viper.SetConfigFile(cfgPath)
//		if err = viper.ReadInConfig(); err == nil {
//			// ensure validateConfig runs properly
//			err = validateConfig(config)
//			if err != nil {
//				return err
//			}
//
//			// marshal the new config
//			out, err := yaml.Marshal(cfg)
//			if err != nil {
//				return err
//			}
//
//			// overwrite the config file
//			err = ioutil.WriteFile(viper.ConfigFileUsed(), out, 0600)
//			if err != nil {
//				return err
//			}
//
//			// set the global variable
//			config = cfg
//		}
//	}
//	return err
//}
