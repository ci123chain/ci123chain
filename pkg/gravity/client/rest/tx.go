package rest

import (
	"fmt"
	"net/http"

	"github.com/ci123chain/ci123chain/pkg/client/context"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	hexUtil "github.com/ethereum/go-ethereum/common/hexutil"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

type valsetConfirmReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	EthAddress string       `json:"eth_address"`
	Nonce      string       `json:"nonce"`
	EthSig     string       `json:"ethSig"`
}

// check the ethereum sig on a particular valset and broadcast a transaction containing
// it if correct. The nonce / block height is used to determine what valset to look up
// locally and verify
func createValsetConfirmHandler(cliCtx context.Context, storeKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req valsetConfirmReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorRes(w, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		res, _, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, req.Nonce), nil, false)
		if err != nil {
			fmt.Printf("could not get valset")
			rest.WriteErrorRes(w, "failed to parse request")
			return
		}
		var valset types.Valset
		cliCtx.Cdc.MustUnmarshalJSON(res, &valset)

		// TODO: fix this, need to fetch the gravityID from params here
		checkpoint := valset.GetCheckpoint("fetch-gravity-id-from-params-please-this-should-panic")

		// the signed message should be the hash of the checkpoint at the given nonce
		ethHash := ethCrypto.Keccak256Hash(checkpoint)

		ethSig, err := hexUtil.Decode(req.EthSig)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		ethPubkey, err := ethCrypto.SigToPub(ethHash.Bytes(), ethSig)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		ethPubkeyBytes := ethCrypto.FromECDSAPub(ethPubkey)

		correct := ethCrypto.VerifySignature(ethPubkeyBytes, ethHash.Bytes(), ethSig)
		if correct == false {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		cosmosAddr := cliCtx.GetFromAddresses()
		msg := types.NewMsgValsetConfirm(valset.Nonce, req.EthAddress, cosmosAddr, req.EthSig)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse request")
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

type bootstrapConfirmReq struct {
	BaseReq               rest.BaseReq           `json:"base_req"`
	Orchestrator          sdk.AccAddress         `json:"orchestrator"`
	EthereumChainID       uint64                 `json:"ethereum_chain_id"`
	BridgeContractAddress string                 `json:"bridge_contract_address"`
	Block                 string                 `json:"block"`
	BridgeValidators      types.BridgeValidators `json:"bridge_validators"`
	GravityID             string                 `json:"gravity_id"`
	StartThreshold        uint64                 `json:"start_threshold"`
}
