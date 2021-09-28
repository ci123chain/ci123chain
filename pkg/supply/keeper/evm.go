package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"math"
	"math/big"
	"strings"
)

///supply/keeper/evm.go

const (
	DefaultDaoABI = `[
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "_holders",
          "type": "address[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_stakes",
          "type": "uint256[]"
        },
        {
          "internalType": "uint64[3]",
          "name": "_votingSettings",
          "type": "uint64[3]"
        },
        {
          "internalType": "uint64",
          "name": "_financePeriod",
          "type": "uint64"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "msgsender",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "thiscontract",
          "type": "address"
        }
      ],
      "name": "ACLAddr",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "acl",
      "outputs": [
        {
          "internalType": "contract ACL",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_sender",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "_role",
          "type": "bytes32"
        },
        {
          "internalType": "uint256[]",
          "name": "_params",
          "type": "uint256[]"
        }
      ],
      "name": "canPerform",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "finance",
      "outputs": [
        {
          "internalType": "contract Finance",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "getInitializationBlock",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "hasInitialized",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "isPetrified",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "contract IACL",
          "name": "_acl",
          "type": "address"
        }
      ],
      "name": "setACL",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "token",
      "outputs": [
        {
          "internalType": "contract MiniMeToken",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "tokenFactory",
      "outputs": [
        {
          "internalType": "contract MiniMeTokenFactory",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "tokenManager",
      "outputs": [
        {
          "internalType": "contract TokenManager",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "vault",
      "outputs": [
        {
          "internalType": "contract Vault",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "vote",
      "outputs": [
        {
          "internalType": "contract Voting",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    }
  ]`

	DefaultTokenManagerABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "msgsender",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "thiscontract",
          "type": "address"
        }
      ],
      "name": "ACLAddr",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "vestingId",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "NewVesting",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "vestingId",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "nonVestedAmount",
          "type": "uint256"
        }
      ],
      "name": "RevokeVesting",
      "type": "event"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "ASSIGN_ROLE",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "BURN_ROLE",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "ISSUE_ROLE",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "MAX_VESTINGS_PER_ADDRESS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "MINT_ROLE",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "REVOKE_VESTINGS_ROLE",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_sender",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "_role",
          "type": "bytes32"
        },
        {
          "internalType": "uint256[]",
          "name": "_params",
          "type": "uint256[]"
        }
      ],
      "name": "canPerform",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "getInitializationBlock",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "hasInitialized",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "isPetrified",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "maxAccountTokens",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "contract IACL",
          "name": "_acl",
          "type": "address"
        }
      ],
      "name": "setACL",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "token",
      "outputs": [
        {
          "internalType": "contract MiniMeToken",
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "vestingsLengths",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "contract MiniMeTokenFactory",
          "name": "_tokenFactory",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "_transferable",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "_maxAccountTokens",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "_symbol",
          "type": "string"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_receiver",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "mint",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "issue",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_receiver",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "assign",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_holder",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "burn",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_receiver",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "internalType": "uint64",
          "name": "_start",
          "type": "uint64"
        },
        {
          "internalType": "uint64",
          "name": "_cliff",
          "type": "uint64"
        },
        {
          "internalType": "uint64",
          "name": "_vested",
          "type": "uint64"
        },
        {
          "internalType": "bool",
          "name": "_revokable",
          "type": "bool"
        }
      ],
      "name": "assignVested",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_holder",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_vestingId",
          "type": "uint256"
        }
      ],
      "name": "revokeVesting",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "_from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "onTransfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "onApprove",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "proxyPayment",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": true,
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "isForwarder",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_vestingId",
          "type": "uint256"
        }
      ],
      "name": "getVesting",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "uint64",
          "name": "start",
          "type": "uint64"
        },
        {
          "internalType": "uint64",
          "name": "cliff",
          "type": "uint64"
        },
        {
          "internalType": "uint64",
          "name": "vesting",
          "type": "uint64"
        },
        {
          "internalType": "bool",
          "name": "revokable",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_holder",
          "type": "address"
        }
      ],
      "name": "spendableBalanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_holder",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_time",
          "type": "uint256"
        }
      ],
      "name": "transferableBalance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "_token",
          "type": "address"
        }
      ],
      "name": "allowRecoverability",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    }
  ]`

	DefaultGas = math.MaxUint64 / 2
	)


// MintCoinsFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) MintCoinsFromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
	//param := []interface{}{metaData.Name, metaData.Symbol, metaData.Symbol, 0, true}
	//denomAddr, err := a.DeployERC20Contract(ctx, owner, param)
	//if err != nil {
	//	return err
	//}
	err := k.Mint(ctx, sdk.HexToAddress(wlkContract), recipientAddr, types.ModuleName, amt)
	return err
}

// TransferFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) TransferFromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
	return k.SendCoinsFromModuleToEVMAccount(ctx, recipientAddr, types.ModuleName, sdk.HexToAddress(wlkContract), amt)
}

func (k Keeper) BuildParams(sender sdk.AccAddress, to *common.Address, payload []byte, gasLimit, nonce uint64) evmtypes.MsgEvmTx {
	return evmtypes.MsgEvmTx{
		From: sender,
		Data: evmtypes.TxData{
			Payload: payload,
			Amount: big.NewInt(0),
			Recipient: to,
			GasLimit: gasLimit,
			AccountNonce: nonce,
		},
	}
}



func (k Keeper) DeployDaoContract(ctx sdk.Context, moduleName string, params interface{}) (sdk.AccAddress, error) {
	ctx.WithIsRootMsg(true)
	sender := k.GetModuleAddress(moduleName)
	abiIns, err := abi.NewABI(DefaultDaoABI)

	if err != nil {
		return sdk.AccAddress{}, err
	}

	bin, err := hex.DecodeString(strings.TrimPrefix(DefaultDaoByteCode, "0x"))

	if err != nil {
		return sdk.AccAddress{}, err
	}

	data, err := abi.Encode(params, abiIns.Constructor.Inputs)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	data = append(bin, data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(sender, nil, data, DefaultGas,  k.getNonce(ctx, sender))
	//msg := k.BuildParams(sender, nil, data, DefaultGas,  0)

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)

	if result != nil && err == nil{
		addStr := result.VMResult().Log
		return sdk.HexToAddress(addStr), err
	}
	return sdk.AccAddress{}, err
}


func (k Keeper) SendCoinsFromEVMAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	moduleName string, wlkContract sdk.AccAddress, amount *big.Int) error {

	to := k.GetModuleAddress(moduleName)

	abiIns, err := abi.NewABI(DefaultDaoABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["transfer"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{to.Address, amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(senderAddr, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, senderAddr))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Transfer action result", "value", result)
	return err
}

func (k Keeper) SendCoinsFromModuleToEVMAccount(ctx sdk.Context, to sdk.AccAddress,
	moduleName string, wlkContract sdk.AccAddress, amount *big.Int) error {

	from := k.GetModuleAddress(moduleName)

	abiIns, err := abi.NewABI(DefaultDaoABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["transfer"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{to.Address, amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(from, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, from))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Transfer action result", "value", result)
	return err
}

func (k Keeper) BurnEVMCoin(ctx sdk.Context, moduleName string, wlkContract, to sdk.AccAddress, amount *big.Int) error {
	from := k.GetModuleAddress(moduleName)

	abiIns, err := abi.NewABI(DefaultTokenManagerABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["burn"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{web3.HexToAddress(to.Address.String()), amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(from, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, from))
	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Burn action result", "value", result)
	return err
}

func (k Keeper) Mint(ctx sdk.Context, contract, to sdk.AccAddress, moduleName string, amount *big.Int) error {
	sender := k.GetModuleAddress(moduleName)
	abiIns, err := abi.NewABI(DefaultTokenManagerABI)
	if err != nil {
		return err
	}

	m, ok := abiIns.Methods["mint"]
	if !ok {
		return fmt.Errorf("invalid method")
	}

	params := []interface{}{web3.HexToAddress(to.Address.String()), amount}

	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(sender, &contract.Address, data, DefaultGas,  k.getNonce(ctx, sender))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Mint action result", "value", result)
	return err
}

func (k Keeper) WRC20DenomValueForFunc(ctx sdk.Context, moduleName string, contract sdk.AccAddress, funcName string) (value interface{}, err error) {
	sender := k.GetModuleAddress(moduleName)
	abiIns, err := abi.NewABI(DefaultDaoABI)
	if err != nil {
		return nil, err
	}
	m, ok := abiIns.Methods[funcName]
	if !ok {
		return nil, fmt.Errorf("invalid method")
	}
	data, err := abi.Encode(nil, m.Inputs)

	data = append(m.ID(), data...)

	msg := k.BuildParams(sender, &contract.Address, data, DefaultGas,  k.getNonce(ctx, sender))

	// for simulator
	ctx.WithIsCheckTx(true)
	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	if err != nil {
		return nil, err
	}
	resData, err := evmtypes.DecodeResultData(result.VMResult().Data)
	if err != nil {
		return nil, err
	}
	if len(resData.Ret) > 0 {
		respInterface, err := abi.Decode(m.Outputs, resData.Ret)
		if err != nil {
			return nil, err
		}
		resp := respInterface.(map[string]interface{})
		v, ok := resp["0"]
		if ok {
			return v, nil
		}
		return nil, err
	}

	return nil, errors.New(fmt.Sprintf("Empty result for func: %s", funcName))
}


func (k Keeper)getNonce(ctx sdk.Context, address sdk.AccAddress) uint64 {
	return k.ak.GetAccount(ctx, address).GetSequence()
}


func readLength(data []byte) (int, error) {
	lengthBig := big.NewInt(0).SetBytes(data[0:32])
	if lengthBig.BitLen() > 63 {
		return 0, fmt.Errorf("length larger than int64: %v", lengthBig.Int64())
	}
	length := int(lengthBig.Uint64())
	if length > len(data) {
		return 0, fmt.Errorf("length insufficient %v require %v", len(data), length)
	}
	return length, nil
}