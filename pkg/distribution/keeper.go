package distribution

import (
	"encoding/hex"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tendermint/tendermint/libs/bech32"
	"strings"
	"sync"

	//"github.com/tanhuiya/ci123chain/pkg/app/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlocker(ak account.AccountKeeper, distr DistrKeeper) types.BeginBlocker{
	return func(ctx types.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock{


		str := Ca(req.Header.ProposerAddress).String()
		res, _ := sdk.ConsAddressFromBech32(str)

		address := "0x" + fmt.Sprintf("%X", []byte(res))
		pAddress := types.HexToAddress(address)


		//var validators []types.AccAddress
		//		//b := req.LastCommitInfo.Votes
		//		//fmt.Println(len(b))
		//		//num := len(b)
		//		//for i := 0; i < num; i++ {
		//		//	p := Ca(b[i].Validator.Address).String()
		//		//	//e := strings.Split(p, "cosmosvalcon")
		//		//	//c := fmt.Sprintf("ox" + e[1])
		//		//	validators[i] = types.HexToAddress(p)
		//		//}

		//fee := types.Coin(uint64(feee))
		//		//proposerAcc := ak.GetAccount(ctx, addr)
		//		//amount := ak.GetBalance(ctx, addr)
		//		//newAmount := amount.SafeAdd(fee)
		//		//err := proposerAcc.SetCoin(newAmount)
		//		//if err != nil {
		//		//	fmt.Print(err)
		//		//}
		//		//ak.SetAccount(ctx, proposerAcc)
		if ctx.BlockHeight() > 1 {
			fee := distr.feeCollectionKeeper.GetCollectedFees(ctx)
			//分配完奖励金之后清空奖金池
			distr.feeCollectionKeeper.ClearCollectedFees(ctx)
			//distr.DistributeRewardsToValidators(ctx, pAddress, fee)

			distr.SetProposerCurrentRewards(ctx, pAddress, fee)
		}

		return abci.ResponseBeginBlock{}
	}
}

type Ca []byte
func (ca Ca) Empty() bool {
	if ca == nil {
		return true
	}
	return false
}

func (ca Ca) Bytes() []byte {
	return ca
}

func (ca Ca) String() string {
	if ca.Empty() {
		return ""
	}
	bech32PrefixConsAddr := GetConfig().GetBech32ConsensusAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, ca.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

type Config struct {
	mtx                 sync.RWMutex
	sealed              bool
	bech32AddressPrefix map[string]string
}


const (
	// AddrLen defines a valid address length
	AddrLen = 20
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32MainPrefix = "cosmos"

	// PrefixAccount is the prefix for account keys
	PrefixAccount = "acc"
	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// PrefixAddress is the prefix for addresses
	PrefixAddress = "addr"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)
var (
	// Initializing an instance of Config
	sdkConfig = &Config{
		sealed: false,
		bech32AddressPrefix: map[string]string{
			"account_addr":   Bech32PrefixAccAddr,
			"validator_addr": Bech32PrefixValAddr,
			"consensus_addr": Bech32PrefixConsAddr,
			"account_pub":    Bech32PrefixAccPub,
			"validator_pub":  Bech32PrefixValPub,
			"consensus_pub":  Bech32PrefixConsPub,
		},
	}
)

func GetConfig() *Config {
	return sdkConfig
}

func (config *Config) GetBech32ConsensusAddrPrefix() string {
	return config.bech32AddressPrefix["consensus_addr"]
}

type HexBytes []byte
func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], []byte(s))
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

// This is the point of Bytes.
func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("Invalid hex string: %s", data)
	}
	bz2, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}