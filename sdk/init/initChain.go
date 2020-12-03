package init

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	acc_type "github.com/ci123chain/ci123chain/pkg/account/types"
	app_module "github.com/ci123chain/ci123chain/pkg/app/module"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	distr "github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/staking"
	stypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	suptypes "github.com/ci123chain/ci123chain/pkg/supply/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
	"regexp"
	"text/template"
	"time"
)

type GenesisFiles struct {
	GenesisBytes			[]byte `json:"genesis_bytes"`
}

type InitFiles struct {
	ConfigBytes				[]byte `json:"config_bytes"`
	PrivValidatorKeyBytes	[]byte `json:"priv_validator_key_bytes"`
	PrivValidatorStateBytes	[]byte `json:"priv_validator_state_bytes"`
	NodeKeyBytes			[]byte `json:"node_key_bytes"`
}

type ChainInfo struct {
	ChainID		string 		`json:"chain_id"`
	GenesisTime time.Time 	`json:"genesis_time"`
}

type ValidatorInfo struct {
	PubKey		string 		`json:"pub_key"`
	Name		string 				`json:"name"`
}

type StakingInfo struct {
	Address 			types.AccAddress	`json:"address"`
	PubKey				string		`json:"pub_key"`
	Tokens				string				`json:"tokens"`
	CommissionInfo  	CommissionInfo 		`json:"commission_info"`
	UpdateTime 			time.Time 			`json:"update_time"`
	MinSelfDelegation   string				`json:"min_self_delegation"`
	Description			stypes.Description  `json:"description"`
}

type CommissionInfo struct {
	Rate        	int64 		`json:"rate"`
	MaxRate     	int64 		`json:"max_rate"`
	MaxChangeRate 	int64		`json:"max_change_rate"`
}

type SupplyInfo struct {
	Amount   string	`json:"amount"`
}

type AccountInfo struct {
	Address  types.AccAddress `json:"address"`
	Amount	 string `json:"amount"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewGenesisFiles(chainInfo ChainInfo,
	validatorInfo []ValidatorInfo, stakingInfo []StakingInfo,
	supplyInfo SupplyInfo, accountInfo []AccountInfo) (*GenesisFiles, error) {

	//todo check infos
	err := checkChainInfo(chainInfo)
	if err != nil {
		return nil, err
	}

	err = checkValidatorInfo(validatorInfo)
	if err != nil {
		return nil, err
	}

	err = checkStakingInfo(stakingInfo)
	if err != nil {
		return nil, err
	}

	err = checkSupplyInfo(supplyInfo)
	if err != nil {
		return nil, err
	}

	err = checkAccountInfo(accountInfo)
	if err != nil {
		return nil, err
	}

	var genesisBytes []byte
	if chainInfo != (ChainInfo{}) {
		//genesis.json
		genesisBytes, err = createGenesis(chainInfo, validatorInfo, stakingInfo, supplyInfo, accountInfo)
		if err != nil {
			return nil, err
		}
	}

	genesisFiles := &GenesisFiles{
		GenesisBytes:            genesisBytes,
	}
	return genesisFiles, nil
}

func NewInitFiles(privKey string, persistentPeers string) (*InitFiles, error){
	//config.toml
	config, err := createConfig(persistentPeers)
	if err != nil {
		return nil, err
	}
	var configTemplate *template.Template
	var buffer bytes.Buffer
	if configTemplate, err = template.New("configFileTemplate").Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}
	configBytes := buffer.Bytes()

	//priv_validator_key(state).json
	privValidatorKeyBytes, privValidatorStateBytes, err := createPrivValidator(privKey)
	if err != nil {
		return nil, err
	}

	//nodeKey.json
	nodeKeyBytes, err := createNodeKey(privKey)
	if err != nil {
		return nil, err
	}

	initFiles := &InitFiles{
		ConfigBytes:             configBytes,
		PrivValidatorKeyBytes:   privValidatorKeyBytes,
		PrivValidatorStateBytes: privValidatorStateBytes,
		NodeKeyBytes: 			 nodeKeyBytes,
	}
	return initFiles, nil
}

func createGenesis(chainInfo ChainInfo, validatorInfo []ValidatorInfo,
	stakingInfo []StakingInfo, supplyInfo SupplyInfo,
	accountInfo []AccountInfo) (genesisBytes []byte, err error) {
	cdc := app_types.MakeCodec()

	var validators []tmtypes.GenesisValidator
	for _, v := range validatorInfo {
		var valKey crypto.PubKey
		pubStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PubKeyAminoName, v.PubKey)
		err = cdc.UnmarshalJSON([]byte(pubStr), &valKey)
		if err != nil {
			return nil, err
		}
		val := app_module.AppGetValidator(valKey, v.Name)
		val.Address =valKey.Address()
		validators = append(validators, val)
	}
	appState := app_module.ModuleBasics.DefaultGenesis(validators)

	err = genesisStakingModule(appState, stakingInfo, cdc)
	if err != nil {
		return nil, err
	}
	genesisDistributionModule(appState, stakingInfo, cdc)

	err = genesisSupplyModule(appState, supplyInfo, cdc)
	if err != nil {
		return nil, err
	}
	err = genesisAccountModule(appState, accountInfo, cdc)
	if err != nil {
		return nil, err
	}

	appStateRaw, err := json.Marshal(appState)
	if err != nil {
		return nil, err
	}
	genDoc := tmtypes.GenesisDoc{
		GenesisTime: 	chainInfo.GenesisTime,
		ChainID: 		chainInfo.ChainID,
		Validators: 	validators,
		AppState:		appStateRaw,
		ConsensusParams: &tmtypes.ConsensusParams{
			Block: tmtypes.DefaultBlockParams(),
			Evidence: tmtypes.DefaultEvidenceParams(),
			Validator: tmtypes.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeSecp256k1}},
		},
	}

	genesisBytes, err = cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return nil, err
	}
	return genesisBytes, nil
}

func createConfig (persistentPeers string) (*cfg.Config, error) {
	c := cfg.DefaultConfig()
	c.Moniker = common.RandStr(8)
	c.Instrumentation.Prometheus = true
	c.Consensus.TimeoutPropose = 5 * time.Second
	c.Consensus.TimeoutCommit = 8 * time.Second
	c.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	c.ProfListenAddress = "localhost:6060"
	c.P2P.RecvRate = 5120000
	c.P2P.SendRate = 5120000
	c.TxIndex.IndexTags = "contract.address,contract.event.data,contract.event.name"
	c.P2P.PersistentPeers = persistentPeers

	err := unmarshalWithViper(viper.GetViper(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func createPrivValidator(privKey string) (privValidatorKey, privValidatorState []byte, err error) {
	validatorKey, err := privStrToPrivKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	cdc := app_types.MakeCodec()
	privValidator := &pvm.FilePV{
		Key:           pvm.FilePVKey{},
		LastSignState: pvm.FilePVLastSignState{},
	}
	privValidator.Key.PrivKey = validatorKey
	privValidator.Key.PubKey = validatorKey.PubKey()
	privValidator.Key.Address = validatorKey.PubKey().Address()

	privValidatorKey, err = cdc.MarshalJSONIndent(privValidator.Key, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	privValidatorState, err = cdc.MarshalJSONIndent(privValidator.LastSignState, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return
}

func createNodeKey(privStr string) (nodeKeyBytes []byte, err error) {
	privKey, err := privStrToPrivKey(privStr)
	if err != nil {
		return nil, err
	}
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}
	cdc := app_types.MakeCodec()
	nodeKeyBytes, err = cdc.MarshalJSON(nodeKey)
	if err != nil {
		return nil, err
	}
	return nodeKeyBytes, nil
}

func genesisStakingModule(appState map[string]json.RawMessage, stakingInfo []StakingInfo, cdc *codec.Codec) error {
	var stakingGenesisState stypes.GenesisState
	var genesisValidators []stypes.Validator
	var delegations []stypes.Delegation
	if _, ok := appState[staking.ModuleName]; !ok{
		stakingGenesisState = stypes.GenesisState{}
	} else {
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
	}

	for _, v := range stakingInfo {
		tokens, ok := types.NewIntFromString(v.Tokens)
		if !ok {
			return errors.New("staking tokens converts to bigInt failed")
		}
		minSelfTokens, ok := types.NewIntFromString(v.MinSelfDelegation)
		if !ok {
			return errors.New("staking minSelfDelegation converts to bigInt failed")
		}
		shares := types.NewDecFromInt(tokens)
		genesisValidators = append(genesisValidators, stypes.Validator{
			OperatorAddress:   v.Address,
			ConsensusKey:      v.PubKey,
			Jailed:            false,
			Status:            1,
			Tokens:            tokens,
			DelegatorShares:   shares,
			Description:       v.Description,
			UnbondingHeight:   -1,
			UnbondingTime:     time.Time{},
			BondedHeight:     0,
			Commission:        stypes.Commission{
				CommissionRates: stypes.CommissionRates{
					Rate:          types.NewDecWithPrec(v.CommissionInfo.Rate, 2),
					MaxRate:       types.NewDecWithPrec(v.CommissionInfo.MaxRate, 2),
					MaxChangeRate: types.NewDecWithPrec(v.CommissionInfo.MaxChangeRate, 2),
				},
				UpdateTime:      v.UpdateTime,
			},
			MinSelfDelegation: minSelfTokens,
		})

		delegations = append(delegations, stypes.NewDelegation(v.Address, v.Address, shares))
	}

	stakingGenesisState.Validators = append(stakingGenesisState.Validators, genesisValidators...)
	stakingGenesisState.Delegations = append(stakingGenesisState.Delegations, delegations...)

	genesisStateBz := cdc.MustMarshalJSON(stakingGenesisState)
	appState[staking.ModuleName] = genesisStateBz
	return nil
}

func genesisSupplyModule(appState map[string]json.RawMessage, supplyInfo SupplyInfo, cdc *codec.Codec) error {
	var supplyGenesisState suptypes.GenesisState
	if _, ok := appState[supply.ModuleName]; !ok{
		supplyGenesisState = suptypes.GenesisState{}
	} else {
		cdc.MustUnmarshalJSON(appState[supply.ModuleName], &supplyGenesisState)
	}

	amount, ok := types.NewIntFromString(supplyInfo.Amount)
	if !ok {
		return errors.New("supply amount converts to bigInt failed")
	}

	supplyGenesisState.Supply = types.NewCoin(amount)
	genesisStateBz := cdc.MustMarshalJSON(supplyGenesisState)
	appState[supply.ModuleName] = genesisStateBz
	return nil
}

func genesisAccountModule(appState map[string]json.RawMessage, accountInfo []AccountInfo, cdc *codec.Codec) error {
	var genesisAccounts acc_type.GenesisAccounts
	if _, ok := appState[account.ModuleName]; !ok {
		genesisAccounts = acc_type.GenesisAccounts{}
	} else {
		cdc.MustUnmarshalJSON(appState[account.ModuleName], &genesisAccounts)
	}

	for _, v := range accountInfo {
		if genesisAccounts.Contains(v.Address) {
			return errors.New(fmt.Sprintf("cannot add account at existing address %v", v.Address))
		}
		amount, ok := types.NewIntFromString(v.Amount)
		if !ok {
			return errors.New("account amount converts to bigInt failed")
		}

		genAcc := account.NewGenesisAccountRaw(v.Address, types.NewCoin(amount))
		if err := genAcc.Validate(); err != nil {
			return err
		}
		genesisAccounts = append(genesisAccounts, genAcc)
	}

	genesisStateBz := cdc.MustMarshalJSON(account.GenesisState(genesisAccounts))
	appState[account.ModuleName] = genesisStateBz
	return nil
}

func genesisDistributionModule(appState map[string]json.RawMessage, stakingInfo []StakingInfo, cdc *codec.Codec) {
	var distributionGenesisState distr.GenesisState
	var outstanddingRewards []distr.ValidatorOutstandingRewardsRecord
	var currentRewards []distr.ValidatorCurrentRewardsRecord
	if _, ok := appState[distr.ModuleName]; !ok {
		distributionGenesisState = distr.GenesisState{}
	}else {
		cdc.MustUnmarshalJSON(appState[distr.ModuleName], &distributionGenesisState)
	}
	for _, v := range stakingInfo {
		outstanddingRewards = append(outstanddingRewards, distr.ValidatorOutstandingRewardsRecord{
			ValidatorAddress:   v.Address,
			OutstandingRewards: types.NewEmptyDecCoin(),
		})
		currentRewards = append(currentRewards, distr.ValidatorCurrentRewardsRecord{
			ValidatorAddress: v.Address,
			Rewards:          distr.ValidatorCurrentRewards{
				Rewards: types.NewEmptyDecCoin(),
				Period:  0,
			},
		})
	}

	distributionGenesisState.ValidatorCurrentRewards = append(distributionGenesisState.ValidatorCurrentRewards, currentRewards...)
	distributionGenesisState.OutstandingRewards = append(distributionGenesisState.OutstandingRewards, outstanddingRewards...)
	distrGenesisStateBz := cdc.MustMarshalJSON(distributionGenesisState)
	appState[distr.ModuleName] = distrGenesisStateBz
}

func privStrToPrivKey(privStr string) (*secp256k1.PrivKeySecp256k1, error) {
	if len(privStr) > 0 {
		//1.match length
		priByt := []byte(privStr)
		length := len(priByt)
		if length != 44 {
			return nil, errors.New(fmt.Sprintf("length of validator key does not match, expected %d, got %d", 44, length))
		}

		//2.regex match
		rule := `=$`
		reg := regexp.MustCompile(rule)
		if !reg.MatchString(privStr) {
			return nil, errors.New("the end of the validator key string should be an equal sign")
		}

		//3.match base64 encoding
		_, err := base64.StdEncoding.DecodeString(privStr)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("privStr cannot be empty")
	}

	var realKey *secp256k1.PrivKeySecp256k1
	privKey := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, privStr)
	cdc := app_types.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(privKey), &realKey)
	if err != nil {
		return nil, err
	}
	return realKey, nil
}

func unmarshalWithViper(vp *viper.Viper, c *cfg.Config) error {
	// you can configure tedermint params via environment variables.
	// TM_PARAMS="consensus.timeout_commit=3000,instrumentation.prometheus=true" ./liamd start
	util.SetEnvToViper(vp, "TM_PARAMS")
	if err := vp.Unmarshal(c); err != nil {
		return err
	}
	return nil
}

//check
func checkChainInfo(chainInfo ChainInfo) error {
	return nil
}

func checkValidatorInfo(validatorInfo []ValidatorInfo) error {
	return nil
}

func checkStakingInfo(stakingInfo []StakingInfo) error {
	return nil
}

func checkSupplyInfo(supplyInfo SupplyInfo) error {
	return nil
}

func checkAccountInfo(accountInfo []AccountInfo) error {
	return nil
}