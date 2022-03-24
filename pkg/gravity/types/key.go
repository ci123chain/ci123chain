package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strings"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	WlKTokenAddress = "0x0000000000000000000000ffff"
)

const (

	// Valsets

	// This retrieves a specific validator set by it's nonce
	// used to compare what's on Ethereum with what's in Cosmos
	// to perform slashing / validation of system consistency
	QueryValsetRequest = "valsetRequest"
	// Gets all the confirmation signatures for a given validator
	// set, used by the relayer to package the validator set and
	// it's signatures into an Ethereum transaction
	QueryValsetConfirmsByNonce = "valsetConfirms"
	// Gets the last N (where N is currently 5) validator sets that
	// have been produced by the chain. Useful to see if any recently
	// signed requests can be submitted.
	QueryLastValsetRequests = "lastValsetRequests"
	// Gets a list of unsigned valsets for a given validators delegate
	// orchestrator address. Up to 100 are sent at a time
	QueryLastPendingValsetRequestByAddr = "lastPendingValsetRequest"

	QueryCurrentValset = "currentValset"
	// TODO remove this, it's not used, getting one confirm at a time
	// is mostly useless
	QueryValsetConfirm = "valsetConfirm"

	// Batches
	// note the current logic here constrains batch throughput to one
	// batch (of any type) per Cosmos block.

	// This retrieves a specific batch by it's nonce and token contract
	// or in the case of a Cosmos originated address it's denom
	QueryBatch = "batch"
	// Get the last unsigned batch (of any denom) for the validators
	// orchestrator to sign
	QueryLastPendingBatchRequestByAddr = "lastPendingBatchRequest"
	// gets the last 100 outgoing batches, regardless of denom, useful
	// for a relayer to see what is available to relay
	QueryLatestTxBatches = "lastBatches"
	// Used by the relayer to package a batch with signatures required
	// to submit to Ethereum
	QueryBatchConfirms = "batchConfirms"
	// Used to query all pending SendToEth transactions and fees available for each
	// token type, a relayer can then estimate their potential profit when requesting
	// a batch
	QueryBatchFees = "batchFees"

	// Logic calls
	// note the current logic here constrains logic call throughput to one
	// call (of any type) per Cosmos block.

	// Token mapping
	// This retrieves the denom which is represented by a given ERC20 contract
	QueryERC20ToDenom = "ERC20ToDenom"
	// This retrieves the ERC20 contract which represents a given denom
	QueryDenomToERC20 = "DenomToERC20"

	// This retrieves the denom which is represented by a given ERC721 contract
	QueryERC721ToDenom = "ERC721ToDenom"
	// This retrieves the ERC721 contract which represents a given denom
	QueryDenomToERC721 = "DenomToERC721"

	// Query pending transactions
	QueryPendingSendToEths = "PendingSendToEth"

	// Query last event nonce
	QueryLastEventNonce = "lastEventNonce"

	// Query last valset confirm nonce
	QueryLastValsetConfirmNonce = "lastValsetConfirmNonce"

	//
	QueryTxId = "txId"
	QueryEventNonce = "eventNonce"

	QueryObservedEventNonce = "observedEventNonce"
)


var (
	// EthAddressKey indexes cosmos validator account addresses
	// i.e. cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	EthAddressKey = []byte{0x1}

	// ValsetRequestKey indexes valset requests by nonce
	ValsetRequestKey = []byte{0x2}

	// ValsetConfirmKey indexes valset confirmations by nonce and the validator account address
	// i.e cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	ValsetConfirmKey = []byte{0x3}

	// OracleClaimKey Claim details by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// A claim is named more intuitively than an Attestation, it is literally
	// a validator making a claim to have seen something happen. Claims are
	// attached to attestations which can be thought of as 'the event' that
	// will eventually be executed.
	OracleClaimKey = []byte{0x4}

	// OracleAttestationKey attestation details by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x5}

	// OutgoingTXPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTXPoolKey = []byte{0x6}


	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0x7}

	// DenomiatorPrefix indexes token contract addresses from ETH on gravity
	DenomiatorPrefix = []byte{0x8}

	// SecondIndexOutgoingTXFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTXFeeKey = []byte{0x9}

	// OutgoingTXBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTXBatchKey = []byte{0xa}

	// OutgoingTXBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTXBatchBlockKey = []byte{0xb}

	// OutgoingTXRequestBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTXRequestBatchKey = []byte{0xc}


	// SecondIndexNonceByClaimKey indexes latest nonce for a given claim type
	SecondIndexNonceByClaimKey = []byte{0xf}


	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xe1}


	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress = []byte{0xe8}

	// KeyOutgoingLogicCall indexes the outgoing logic calls
	KeyOutgoingLogicCall = []byte{0xde}

	// KeyOutgoingLogicConfirm indexes the outgoing logic confirms
	KeyOutgoingLogicConfirm = []byte{0xae}

	// save gravity list
	GravityListKey = []byte{0xaf}

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xf1}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xf2}

	// LastValsetConfirmNonceKey indexes the latest valset confirm nonce
	LastValsetConfirmNonceKey = []byte{0xf3}


	// LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	LastObservedEthereumBlockHeightKey = []byte{0xf4}

	// WlkToEthKey prefixes the index of wlk asset to eth ERC20s
	WlkToEthKey = []byte{0xf5}

	// EthToWlkKey prefixes the index of eth originated assets ERC20s to wlk
	EthToWlkKey = []byte{0xf6}

	// LastSlashedValsetNonce indexes the latest slashed valset nonce
	LastSlashedValsetNonce = []byte{0xf7}

	// LatestValsetNonce indexes the latest valset nonce
	LatestValsetNonce = []byte{0xf8}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0xf9}

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0xfa}

	ContractMetaDataKey = []byte{0xfb}

	TxIdKey = []byte{0xfc}

	EventNonceKey = []byte{0xfd}

	// WlkToEthKey prefixes the index of wlk asset to eth ERC20s
	WRC721ToEth721Key = []byte{0xfe}

	// EthToWlkKey prefixes the index of eth originated assets ERC20s to wlk
	ERC721ToWRC721Key = []byte{0xff}


)

// GetOrchestratorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return append(KeyOrchestratorAddress, orc.Bytes()...)
}

// GetEthAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthAddressKey(validator sdk.AccAddress) []byte {
	return append(EthAddressKey, validator.Bytes()...)
}

// GetValsetKey returns the following key format
// prefix    nonce
// [0x0][0 0 0 0 0 0 0 1]
func GetValsetKey(nonce uint64) []byte {
	return append(ValsetRequestKey, UInt64Bytes(nonce)...)
}

func GetGravityKey(gid string) []byte {
	return append(GravityListKey, []byte(gid)...)
}

// GetValsetConfirmKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetValsetConfirmKey(nonce uint64, validator sdk.AccAddress) []byte {
	return append(ValsetConfirmKey, append(UInt64Bytes(nonce), validator.Bytes()...)...)
}

// GetClaimKey returns the following key format
// prefix type               cosmos-validator-address                       nonce                             attestation-details-hash
// [0x0][0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// The Claim hash identifies a unique event, for example it would have a event nonce, a sender and a receiver. Or an event nonce and a batch nonce. But
// the Claim is stored indexed with the claimer key to make sure that it is unique.
func GetClaimKey(details EthereumClaim) []byte {
	var detailsHash []byte
	if details != nil {
		detailsHash = details.ClaimHash()
	} else {
		panic("No claim without details!")
	}
	claimTypeLen := len([]byte{byte(details.GetType())})
	nonceBz := UInt64Bytes(details.GetEventNonce())
	key := make([]byte, len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz)+len(detailsHash))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], []byte{byte(details.GetType())})
	// TODO this is the delegate address, should be stored by the valaddress
	copy(key[len(OracleClaimKey)+claimTypeLen:], details.GetClaimer().Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen:], nonceBz)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz):], detailsHash)
	return key
}

// GetAttestationKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(UInt64Bytes(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], UInt64Bytes(eventNonce))
	copy(key[len(OracleAttestationKey)+len(UInt64Bytes(0)):], claimHash)
	return key
}

// GetAttestationKeyWithHash returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKeyWithHash(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(UInt64Bytes(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], UInt64Bytes(eventNonce))
	copy(key[len(OracleAttestationKey)+len(UInt64Bytes(0)):], claimHash)
	return key
}

// GetOutgoingTxPoolKey returns the following key format
// prefix     id
// [0x6][0 0 0 0 0 0 0 1]
func GetOutgoingTxPoolKey(id uint64) []byte {
	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
}

// GetOutgoingTxBatchKey returns the following key format
// prefix     nonce                     eth-contract-address
// [0xa][0 0 0 0 0 0 0 1][0xc783df8a850f42e7F7e57013759C285caa701eB6]
func GetOutgoingTxBatchKey(tokenContract string, nonce uint64) []byte {
	tokenContract = strings.ToLower(tokenContract)
	return append(append(OutgoingTXBatchKey, []byte(tokenContract)...), UInt64Bytes(nonce)...)
}

func GetOutgoingTxRequestBatchKey(tokenContract string, nonce uint64) []byte {
	tokenContract = strings.ToLower(tokenContract)
	return append(append(OutgoingTXRequestBatchKey, []byte(tokenContract)...), UInt64Bytes(nonce)...)
}

// GetOutgoingTxBatchBlockKey returns the following key format
// prefix     blockheight
// [0xb][0 0 0 0 2 1 4 3]
func GetOutgoingTxBatchBlockKey(block uint64) []byte {
	return append(OutgoingTXBatchBlockKey, UInt64Bytes(block)...)
}

// GetBatchConfirmKey returns the following key format
// prefix           eth-contract-address                BatchNonce                       Validator-address
// [0xe1][0xc783df8a850f42e7F7e57013759C285caa701eB6][0 0 0 0 0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// TODO this should be a sdk.AccAddress
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, validator sdk.AccAddress) []byte {
	tokenContract = strings.ToLower(tokenContract)
	a := append(UInt64Bytes(batchNonce), validator.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetFeeSecondIndexKey returns the following key format
// prefix            eth-contract-address            fee_amount
// [0x9][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000]
func GetFeeSecondIndexKey(fee ERC20Token) []byte {
	r := make([]byte, 1+ETHContractAddressLen+32)
	// sdkInts have a size limit of 255 bits or 32 bytes
	// therefore this will never panic and is always safe
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	// TODO this won't ever work fix it
	copy(r[0:], SecondIndexOutgoingTXFeeKey)
	copy(r[len(SecondIndexOutgoingTXFeeKey):], []byte(fee.Contract))
	copy(r[len(SecondIndexOutgoingTXFeeKey)+len(fee.Contract):], amount)
	return r
}

// GetLastEventNonceByValidatorKey indexes lateset event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetLastEventNonceByValidatorKey(validator sdk.AccAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

func GetWlKToEthKey(denom string) []byte {
	denom = strings.ToLower(denom)
	return append(WlkToEthKey, []byte(denom)...)
}

func GetEthToWlkKey(erc20 string) []byte {
	erc20 = strings.ToLower(erc20)
	return append(EthToWlkKey, []byte(erc20)...)
}

func GetWRC721ToERC721Key(wrc721 string) []byte {
	wrc721 = strings.ToLower(wrc721)
	return append(WRC721ToEth721Key, []byte(wrc721)...)
}

func GetERC721ToWRC721Key(erc721 string) []byte {
	erc721 = strings.ToLower(erc721)
	return append(ERC721ToWRC721Key, []byte(erc721)...)
}

func GetContractMetaDataKey(contract string) []byte {
	return append(ContractMetaDataKey, []byte(contract)...)
}

//func GetOutgoingLogicCallKey(invalidationId []byte, invalidationNonce uint64) []byte {
//	a := append(KeyOutgoingLogicCall, invalidationId...)
//	return append(a, UInt64Bytes(invalidationNonce)...)
//}
//
//func GetLogicConfirmKey(invalidationId []byte, invalidationNonce uint64, validator sdk.AccAddress) []byte {
//	interm := append(KeyOutgoingLogicConfirm, invalidationId...)
//	interm = append(interm, UInt64Bytes(invalidationNonce)...)
//	return append(interm, validator.Bytes()...)
//}

func GetTxIdKey(txId uint64) []byte {
	return append(TxIdKey, sdk.Uint64ToBigEndian(txId)...)
}

func GetEventNonceKey(eventNonce uint64) []byte {
	return append(EventNonceKey, sdk.Uint64ToBigEndian(eventNonce)...)
}