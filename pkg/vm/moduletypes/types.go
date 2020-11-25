package moduletypes

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	RouteKey = "vm"
	ModuleName = "vm"
	StoreKey = "vm"
	DefaultCodespace = "vm"
)

// go to ../keeper/keeper.go
type KeeperI interface {
	BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock)

	EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate

	Upload(ctx sdk.Context, wasmCode []byte, creator sdk.AccAddress) (codeHash []byte, err error)

	Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args utils.CallData, name, version, author, email, describe string, genesisContractAddress sdk.AccAddress) (sdk.AccAddress, error)

	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args utils.CallData) (sdk.Result, error)

	SetBalance(ctx sdk.Context, addr ethcmn.Address, amount *big.Int)

	AddBalance(ctx sdk.Context, addr ethcmn.Address, amount *big.Int)

	SubBalance(ctx sdk.Context, addr ethcmn.Address, amount *big.Int)

	SetNonce(ctx sdk.Context, addr ethcmn.Address, nonce uint64)

	SetState(ctx sdk.Context, addr ethcmn.Address, key, value ethcmn.Hash)

	SetCode(ctx sdk.Context, addr ethcmn.Address, code []byte)

	SetLogs(ctx sdk.Context, hash ethcmn.Hash, logs []*ethtypes.Log) error

	Finalise(ctx sdk.Context, deleteEmptyObjects bool) error

	Commit(ctx sdk.Context, deleteEmptyObjects bool) (root ethcmn.Hash, err error)

	SetChainConfig(ctx sdk.Context, config evmtypes.ChainConfig)

	SetParams(ctx sdk.Context, params evmtypes.Params)
}

