package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/types"
)

// DefaultParamspace defines the default auth module parameter subspace
const (
	// todo: implement oracle constants as params
	DefaultParamspace = ModuleName
	AttestationPeriod = 24 * time.Hour // TODO: value????
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeyGravityID stores the gravity id
	ParamsStoreKeyGravityID = []byte("GravityID")

	// ParamsStoreKeyContractHash stores the contract hash
	ParamsStoreKeyContractHash = []byte("ContractHash")

	// ParamsStoreKeyStartThreshold stores the start threshold
	ParamsStoreKeyStartThreshold = []byte("StartThreshold")

	// ParamsStoreKeyBridgeContractAddress stores the contract address
	ParamsStoreKeyBridgeContractAddress = []byte("BridgeContractAddress")

	// ParamsStoreKeyBridgeContractChainID stores the bridge chain id
	ParamsStoreKeyBridgeContractChainID = []byte("BridgeChainID")

	// ParamsStoreKeySignedValsetsWindow stores the signed blocks window
	ParamsStoreKeySignedValsetsWindow = []byte("SignedValsetsWindow")

	// ParamsStoreKeySignedBatchesWindow stores the signed blocks window
	ParamsStoreKeySignedBatchesWindow = []byte("SignedBatchesWindow")

	// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
	ParamsStoreKeySignedClaimsWindow = []byte("SignedClaimsWindow")

	// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
	ParamsStoreKeyTargetBatchTimeout = []byte("TargetBatchTimeout")

	// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
	ParamsStoreKeyAverageBlockTime = []byte("AverageBlockTime")

	// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
	ParamsStoreKeyAverageEthereumBlockTime = []byte("AverageEthereumBlockTime")

	// ParamsStoreSlashFractionValset stores the slash fraction valset
	ParamsStoreSlashFractionValset = []byte("SlashFractionValset")

	// ParamsStoreSlashFractionBatch stores the slash fraction Batch
	ParamsStoreSlashFractionBatch = []byte("SlashFractionBatch")

	// ParamsStoreSlashFractionClaim stores the slash fraction Claim
	ParamsStoreSlashFractionClaim = []byte("SlashFractionClaim")

	// ParamsStoreSlashFractionConflictingClaim stores the slash fraction ConflictingClaim
	ParamsStoreSlashFractionConflictingClaim = []byte("SlashFractionConflictingClaim")

	//  ParamStoreUnbondSlashingValsetsWindow stores unbond slashing valset window
	ParamStoreUnbondSlashingValsetsWindow = []byte("UnbondSlashingValsetsWindow")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// GenesisState struct
type GenesisState struct {
	Params             *Params                      `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	LastObservedNonce  uint64                       `protobuf:"varint,2,opt,name=last_observed_nonce,json=lastObservedNonce,proto3" json:"last_observed_nonce,omitempty"`
	Valsets            []*Valset                    `protobuf:"bytes,3,rep,name=valsets,proto3" json:"valsets,omitempty"`
	ValsetConfirms     []*MsgValsetConfirm          `protobuf:"bytes,4,rep,name=valset_confirms,json=valsetConfirms,proto3" json:"valset_confirms,omitempty"`
	Batches            []*OutgoingTxBatch           `protobuf:"bytes,5,rep,name=batches,proto3" json:"batches,omitempty"`
	BatchConfirms      []MsgConfirmBatch            `protobuf:"bytes,6,rep,name=batch_confirms,json=batchConfirms,proto3" json:"batch_confirms"`
	LogicCalls         []*OutgoingLogicCall         `protobuf:"bytes,7,rep,name=logic_calls,json=logicCalls,proto3" json:"logic_calls,omitempty"`
	LogicCallConfirms  []MsgConfirmLogicCall        `protobuf:"bytes,8,rep,name=logic_call_confirms,json=logicCallConfirms,proto3" json:"logic_call_confirms"`
	Attestations       []Attestation                `protobuf:"bytes,9,rep,name=attestations,proto3" json:"attestations"`
	DelegateKeys       []*MsgSetOrchestratorAddress `protobuf:"bytes,10,rep,name=delegate_keys,json=delegateKeys,proto3" json:"delegate_keys,omitempty"`
	Erc20ToDenoms      []*ERC20ToDenom              `protobuf:"bytes,11,rep,name=erc20_to_denoms,json=erc20ToDenoms,proto3" json:"erc20_to_denoms,omitempty"`
	UnbatchedTransfers []*OutgoingTransferTx        `protobuf:"bytes,12,rep,name=unbatched_transfers,json=unbatchedTransfers,proto3" json:"unbatched_transfers,omitempty"`
}

type Params struct {
	GravityId                     string                                 `protobuf:"bytes,1,opt,name=gravity_id,json=gravityId,proto3" json:"gravity_id,omitempty"`
	ContractSourceHash            string                                 `protobuf:"bytes,2,opt,name=contract_source_hash,json=contractSourceHash,proto3" json:"contract_source_hash,omitempty"`
	BridgeEthereumAddress         string                                 `protobuf:"bytes,4,opt,name=bridge_ethereum_address,json=bridgeEthereumAddress,proto3" json:"bridge_ethereum_address,omitempty"`
	BridgeChainId                 uint64                                 `protobuf:"varint,5,opt,name=bridge_chain_id,json=bridgeChainId,proto3" json:"bridge_chain_id,omitempty"`
	SignedValsetsWindow           uint64                                 `protobuf:"varint,6,opt,name=signed_valsets_window,json=signedValsetsWindow,proto3" json:"signed_valsets_window,omitempty"`
	SignedBatchesWindow           uint64                                 `protobuf:"varint,7,opt,name=signed_batches_window,json=signedBatchesWindow,proto3" json:"signed_batches_window,omitempty"`
	SignedClaimsWindow            uint64                                 `protobuf:"varint,8,opt,name=signed_claims_window,json=signedClaimsWindow,proto3" json:"signed_claims_window,omitempty"`
	TargetBatchTimeout            uint64                                 `protobuf:"varint,10,opt,name=target_batch_timeout,json=targetBatchTimeout,proto3" json:"target_batch_timeout,omitempty"`
	AverageBlockTime              uint64                                 `protobuf:"varint,11,opt,name=average_block_time,json=averageBlockTime,proto3" json:"average_block_time,omitempty"`
	AverageEthereumBlockTime      uint64                                 `protobuf:"varint,12,opt,name=average_ethereum_block_time,json=averageEthereumBlockTime,proto3" json:"average_ethereum_block_time,omitempty"`
	SlashFractionValset           sdk.Dec `protobuf:"bytes,13,opt,name=slash_fraction_valset,json=slashFractionValset,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_valset"`
	SlashFractionBatch            sdk.Dec `protobuf:"bytes,14,opt,name=slash_fraction_batch,json=slashFractionBatch,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_batch"`
	SlashFractionClaim            sdk.Dec `protobuf:"bytes,15,opt,name=slash_fraction_claim,json=slashFractionClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_claim"`
	SlashFractionConflictingClaim sdk.Dec `protobuf:"bytes,16,opt,name=slash_fraction_conflicting_claim,json=slashFractionConflictingClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_conflicting_claim"`
	UnbondSlashingValsetsWindow   uint64                                 `protobuf:"varint,17,opt,name=unbond_slashing_valsets_window,json=unbondSlashingValsetsWindow,proto3" json:"unbond_slashing_valsets_window,omitempty"`
}

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// TODO: set some better defaults here
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// DefaultParams returns a copy of the default params
func DefaultParams() *Params {
	return &Params{
		GravityId:                     "defaultgravityid",
		SignedValsetsWindow:           10000,
		SignedBatchesWindow:           10000,
		SignedClaimsWindow:            10000,
		TargetBatchTimeout:            43200000,
		AverageBlockTime:              5000,
		AverageEthereumBlockTime:      15000,
		SlashFractionValset:           sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBatch:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionClaim:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionConflictingClaim: sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		UnbondSlashingValsetsWindow:   10000,
	}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validateGravityID(p.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateContractHash(p.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateBridgeContractAddress(p.BridgeEthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	if err := validateTargetBatchTimeout(p.TargetBatchTimeout); err != nil {
		return sdkerrors.Wrap(err, "Batch timeout")
	}
	if err := validateAverageBlockTime(p.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Block time")
	}
	if err := validateAverageEthereumBlockTime(p.AverageEthereumBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Ethereum block time")
	}
	if err := validateSignedValsetsWindow(p.SignedValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedBatchesWindow(p.SignedBatchesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedClaimsWindow(p.SignedClaimsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionValset(p.SlashFractionValset); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionBatch(p.SlashFractionBatch); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionClaim(p.SlashFractionClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionConflictingClaim(p.SlashFractionConflictingClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateUnbondSlashingValsetsWindow(p.UnbondSlashingValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond Slashing valset window")
	}

	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &p.GravityId, validateGravityID),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeEthereumAddress, validateBridgeContractAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedValsetsWindow, &p.SignedValsetsWindow, validateSignedValsetsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedBatchesWindow, &p.SignedBatchesWindow, validateSignedBatchesWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedClaimsWindow, &p.SignedClaimsWindow, validateSignedClaimsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetBatchTimeout, &p.TargetBatchTimeout, validateTargetBatchTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &p.AverageEthereumBlockTime, validateAverageEthereumBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionValset, &p.SlashFractionValset, validateSlashFractionValset),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionBatch, &p.SlashFractionBatch, validateSlashFractionBatch),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionClaim, &p.SlashFractionClaim, validateSlashFractionClaim),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingClaim, &p.SlashFractionConflictingClaim, validateSlashFractionConflictingClaim),
		paramtypes.NewParamSetPair(ParamStoreUnbondSlashingValsetsWindow, &p.UnbondSlashingValsetsWindow, validateUnbondSlashingValsetsWindow),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := GravityCodec.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := GravityCodec.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := strToFixByteArray(v); err != nil {
		return err
	}
	return nil
}

func validateContractHash(i interface{}) error {
	// TODO: should we validate that the input here is a properly formatted
	// SHA256 (or other) hash?
	if _, ok := i.(string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBridgeChainID(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateTargetBatchTimeout(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	return nil
}

func validateAverageBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Cosmos block time, too short for latency limitations")
	}
	return nil
}

func validateAverageEthereumBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Ethereum block time, too short for latency limitations")
	}
	return nil
}

func validateBridgeContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := ValidateEthAddress(v); err != nil {
		// TODO: ensure that empty addresses are valid in params
		if !strings.Contains(err.Error(), "empty") {
			return err
		}
	}
	return nil
}

func validateSignedValsetsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUnbondSlashingValsetsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionValset(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedBatchesWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedClaimsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBatch(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionClaim(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionConflictingClaim(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
