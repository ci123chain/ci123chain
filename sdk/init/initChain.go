package init

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	acc_type "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	distr "github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/staking"
	stypes"github.com/ci123chain/ci123chain/pkg/staking/types"
	suptypes"github.com/ci123chain/ci123chain/pkg/supply/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
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

type InitFiles struct {
	ConfigBytes				[]byte `json:"config_bytes"`
	PrivValidatorKeyBytes	[]byte `json:"priv_validator_key_bytes"`
	PrivValidatorStateBytes	[]byte `json:"priv_validator_state_bytes"`
	NodeKeyBytes			[]byte `json:"node_key_bytes"`
	GenesisBytes			[]byte `json:"genesis_bytes"`
}

type ChainInfo struct {
	ChainID		string 		`json:"chain_id"`
	GenesisTime time.Time 	`json:"genesis_time"`
}

type ValidatorInfo struct {
	PubKey		crypto.PubKey 		`json:"pub_key"`
	Name		string 				`json:"name"`
}

type StakingInfo struct {
	Address 			types.AccAddress	`json:"address"`
	PubKey				crypto.PubKey		`json:"pub_key"`
	Tokens				int64 				`json:"tokens"`
	CommissionInfo  	CommissionInfo 		`json:"commission_info"`
	UpdateTime 			time.Time 			`json:"update_time"`
}

type CommissionInfo struct {
	Rate        	int64 		`json:"rate"`
	MaxRate     	int64 		`json:"max_rate"`
	MaxChangeRate 	int64		`json:"max_change_rate"`
}

type SupplyInfo struct {
	Amount   uint64	`json:"amount"`
}

type AccountInfo struct {
	Address  types.AccAddress `json:"address"`
	Amount	 uint64 `json:"amount"`
}

func NewInitChainFiles(chainInfo ChainInfo,
	validatorInfo ValidatorInfo, stakingInfo StakingInfo,
	supplyInfo SupplyInfo, accountInfo AccountInfo,
	privKey, persistentPeers string) (*InitFiles, error) {

	//todo check infos


	//genesis.json
	genesisBytes, err := createGenesis(chainInfo, validatorInfo, stakingInfo, supplyInfo, accountInfo, privKey)
	if err != nil {
		return nil, err
	}

	//config.toml
	ctx := app.NewDefaultContext()
	err = createConfig(ctx, persistentPeers)
	if err != nil {
		return nil, err
	}
	var configTemplate *template.Template
	var buffer bytes.Buffer
	if configTemplate, err = template.New("configFileTemplate").Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
	if err := configTemplate.Execute(&buffer, ctx.Config); err != nil {
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
		GenesisBytes:            genesisBytes,
	}
	return initFiles, nil
}

func createGenesis(chainInfo ChainInfo, validatorInfo ValidatorInfo,
	stakingInfo StakingInfo, supplyInfo SupplyInfo,
	accountInfo AccountInfo, privKey string) (genesisBytes []byte, err error) {
	validatorKey, err := privStrToPrivKey(privKey)
	if err != nil {
		return nil, err
	}

	cdc := app.MakeCodec()
	val := app.AppGetValidator(validatorKey.PubKey(), validatorInfo.Name)
	val.Address = validatorKey.PubKey().Address()
	validators := []tmtypes.GenesisValidator{val}
	appState := app.ModuleBasics.DefaultGenesis(validators)

	genesisStakingModule(appState, *validatorKey, stakingInfo, cdc)
	genesisSupplyModule(appState, supplyInfo, cdc)
	err = genesisAccountModule(appState, accountInfo, cdc)
	if err != nil {
		return nil, err
	}
	genesisDistributionModule(appState, stakingInfo, cdc)

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

func createConfig(ctx *app.Context, persistentPeers string) error {
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
		return err
	}
	ctx.Config = c
	return nil
}

func createPrivValidator(privKey string) (privValidatorKey, privValidatorState []byte, err error) {
	validatorKey, err := privStrToPrivKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	cdc := app.MakeCodec()
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
	cdc := app.MakeCodec()
	nodeKeyBytes, err = cdc.MarshalJSON(nodeKey)
	if err != nil {
		return nil, err
	}
	return nodeKeyBytes, nil
}

func genesisStakingModule(appState map[string]json.RawMessage, validatorKey secp256k1.PrivKeySecp256k1, stakingInfo StakingInfo, cdc *codec.Codec)  {
	var stakingGenesisState stypes.GenesisState
	var genesisValidator stypes.Validator
	if _, ok := appState[staking.ModuleName]; !ok{
		stakingGenesisState = stypes.GenesisState{}
	} else {
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
	}

	genesisValidator = stypes.Validator{
		OperatorAddress:   stakingInfo.Address,
		ConsensusKey:      hex.EncodeToString(cdc.MustMarshalJSON(validatorKey.PubKey())),
		Jailed:            false,
		Status:            1,
		Tokens:            types.NewInt(stakingInfo.Tokens),
		DelegatorShares:   types.NewDec(stakingInfo.Tokens),
		Description:       stypes.Description{},
		UnbondingHeight:   0,
		UnbondingTime:     time.Time{},
		Commission:        stypes.Commission{
			CommissionRates: stypes.CommissionRates{
				Rate:          types.NewDecWithPrec(stakingInfo.CommissionInfo.Rate, 2),
				MaxRate:       types.NewDecWithPrec(stakingInfo.CommissionInfo.MaxRate, 2),
				MaxChangeRate: types.NewDecWithPrec(stakingInfo.CommissionInfo.MaxChangeRate, 2),
			},
			UpdateTime:      stakingInfo.UpdateTime,
		},
		MinSelfDelegation: types.NewInt(stakingInfo.Tokens),
	}

	delegation := stypes.NewDelegation(stakingInfo.Address, stakingInfo.Address, types.NewDec(stakingInfo.Tokens))

	stakingGenesisState.Validators = append(stakingGenesisState.Validators, genesisValidator)
	stakingGenesisState.Delegations = append(stakingGenesisState.Delegations, delegation)

	genesisStateBz := cdc.MustMarshalJSON(stakingGenesisState)
	appState[staking.ModuleName] = genesisStateBz
	return
}

func genesisSupplyModule(appState map[string]json.RawMessage, supplyInfo SupplyInfo, cdc *codec.Codec) {
	var supplyGenesisState suptypes.GenesisState
	if _, ok := appState[supply.ModuleName]; !ok{
		supplyGenesisState = suptypes.GenesisState{}
	} else {
		cdc.MustUnmarshalJSON(appState[supply.ModuleName], &supplyGenesisState)
	}

	supplyGenesisState.Supply = types.NewUInt64Coin(supplyInfo.Amount)
	genesisStateBz := cdc.MustMarshalJSON(supplyGenesisState.Supply)
	appState[supply.ModuleName] = genesisStateBz
	return
}

func genesisAccountModule(appState map[string]json.RawMessage, accountInfo AccountInfo, cdc *codec.Codec) error{
	var genesisAccounts acc_type.GenesisAccounts
	if _, ok := appState[account.ModuleName]; !ok {
		genesisAccounts = acc_type.GenesisAccounts{}
	} else {
		cdc.MustUnmarshalJSON(appState[account.ModuleName], &genesisAccounts)
	}
	if genesisAccounts.Contains(accountInfo.Address) {
		_ = fmt.Errorf("cannot add account at existing address %v", accountInfo.Address)
	}
	genAcc := account.NewGenesisAccountRaw(accountInfo.Address, types.NewUInt64Coin(accountInfo.Amount))
	if err := genAcc.Validate(); err != nil {
		return err
	}
	genesisAccounts = append(genesisAccounts, genAcc)
	genesisStateBz := cdc.MustMarshalJSON(account.GenesisState(genesisAccounts))
	appState[account.ModuleName] = genesisStateBz
	return nil
}

func genesisDistributionModule(appState map[string]json.RawMessage, stakingInfo StakingInfo, cdc *codec.Codec) {
	var distributionGenesisState distr.GenesisState
	if _, ok := appState[distr.ModuleName]; !ok {
		distributionGenesisState = distr.GenesisState{}
	}else {
		cdc.MustUnmarshalJSON(appState[distr.ModuleName], &distributionGenesisState)
	}
	outstanddingReward := distr.ValidatorOutstandingRewardsRecord{
		ValidatorAddress:   stakingInfo.Address,
		OutstandingRewards: types.NewEmptyDecCoin(),
	}
	currentReward := distr.ValidatorCurrentRewardsRecord{
		ValidatorAddress: stakingInfo.Address,
		Rewards:          distr.ValidatorCurrentRewards{
			Rewards: types.NewEmptyDecCoin(),
			Period:  0,
		},
	}
	distributionGenesisState.ValidatorCurrentRewards = append(distributionGenesisState.ValidatorCurrentRewards, currentReward)
	distributionGenesisState.OutstandingRewards = append(distributionGenesisState.OutstandingRewards, outstanddingReward)
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
	cdc := app.MakeCodec()
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
