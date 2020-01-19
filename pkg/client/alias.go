package client

import "github.com/tanhuiya/ci123chain/pkg/client/types"

var (
	ErrNewClientCtx		= types.ErrNewClientCtx
	ErrGetInputAddr		= types.ErrGetInputAddrCtx
	ErrParseAddr		= types.ErrParseAddr
	ErrNoAddr       	= types.ErrNoAddr
	ErrGetPassPhrase	= types.ErrGetPassPhrase
	ErrGetSignData		= types.ErrGetSignData
	ErrBroadcast		= types.ErrBroadcast
	ErrGetCheckPassword	= types.ErrGetCheckPassword
	ErrGetPassword		= types.ErrGetPassword
	ErrPhrasesNotMatch	= types.ErrPhrasesNotMatch
	ErrNode				= types.ErrNode
)


type BlockIDParts struct {
	//
	Total     string    `json:"total"`
	Hash      string    `json:"hash"`
}

type HeaderVersion struct {
	//
	Block    string    `json:"block"`
	App      string     `json:"app"`
}


type BlockId struct {
	//
	Hash    string  `json:"hash"`
	Parts   BlockIDParts  `json:"parts"`
}

type BlockMeta  struct {
	//
	Blockid    BlockId    `json:"block_id"`
	Header     BlockHeader  `json:"header"`
}

type BlockHeader struct {
	Version HeaderVersion  `json:"version"`
	Chainid string         `json:"chain_id"`
	Height  string         `json:"height"`
	Time    string         `json:"time"`
	Numtxs  string         `json:"num_txs"`
	Totaltxs string        `json:"total_txs"`
	Lastblockid  BlockId    `json:"last_block_id"`
	Lastcommithash string   `json:"last_commit_hash"`
	Datahash     string     `json:"data_hash"`
	Validatorshash  string   `json:"validators_hash"`
	Nextvalidatorshash string  `json:"next_validators_hash"`
	Consensushash    string    `json:"consensus_hash"`
	Apphash        string     `json:"app_hash"`
	Lastresulthash  string    `json:"last_result_hash"`
	Evidencehash    string     `json:"evidence_hash"`
	Proposetaddress  string     `json:"proposer_address"`

}

type PreCommitsInfo struct {
	Type  int64    `json:"type"`
	Height  string   `json:"height"`
	Round   string   `json:"round"`
	Blockid  BlockId  `json:"block_id"`
	Timestamp  string   `json:"timestamp"`
	Validatoraddress string `json:"validator_address"`
	Validatorindex string `json:"validator_index"`
	Signature   string     `json:"signature"`
}

type LastCommit struct {
	Blockid   BlockId   `json:"block_id"`
	Precommits  []PreCommitsInfo  `json:"precommits"`
}

type BlockData struct {
	Txs     []byte   `json:"tx"`
}

type BlockkEvidence struct {
	Evidence   []byte    `json:"evidence"`
}

type BlockInfo struct {
	Header  BlockHeader    `json:"header"`
	Data    BlockData    `json:"data"`
	Evidence  BlockkEvidence   `json:"evidence"`
	Lastcommit  LastCommit  `json:"last_commit"`
	//Precommits   []PreCommitsInfo   `json:"precommits"`
}

type  BlockResult struct {
	Blockmeta    BlockMeta  `json:"block_meta"`
	Block       BlockInfo    `json:"block"`
}


type BlockInformation struct {
	Jsonrpc    string    `json:"jsonrpc"`
	ID         string     `json:"id"`
	Result     BlockResult `json:"result"`
}