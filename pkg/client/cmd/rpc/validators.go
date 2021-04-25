package rpc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"net/http"
	"strconv"
)

// Latest Validator Set REST handler
func LatestValidatorSetRequestHandlerFn(clientCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse pagination parameters")
			return
		}

		output, err := GetValidators(clientCtx, nil, &page, &limit)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}


// Validator output in bech32 format
type ValidatorOutput struct {
	Address          sdk.AccAddress    `json:"address"`
	PubKey           crypto.PubKey `json:"pub_key"`
	ProposerPriority int64              `json:"proposer_priority"`
	VotingPower      int64              `json:"voting_power"`
}

// Validators at a certain height output in bech32 format
type ResultValidatorsOutput struct {
	BlockHeight int64             `json:"block_height"`
	Validators  []ValidatorOutput `json:"validators"`
}

// GetValidators from client
func GetValidators(clientCtx context.Context, height *int64, page, limit *int) (ResultValidatorsOutput, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return ResultValidatorsOutput{}, err
	}

	validatorsRes, err := node.Validators(clientCtx.Context(), height, page, limit)
	if err != nil {
		return ResultValidatorsOutput{}, err
	}

	outputValidatorsRes := ResultValidatorsOutput{
		BlockHeight: validatorsRes.BlockHeight,
		Validators:  make([]ValidatorOutput, len(validatorsRes.Validators)),
	}

	for i := 0; i < len(validatorsRes.Validators); i++ {
		outputValidatorsRes.Validators[i], err = validatorOutput(validatorsRes.Validators[i])
		if err != nil {
			return ResultValidatorsOutput{}, err
		}
	}

	return outputValidatorsRes, nil
}

func validatorOutput(validator *tmtypes.Validator) (ValidatorOutput, error) {
	return ValidatorOutput{
		Address:          sdk.ToAccAddress(validator.Address.Bytes()),
		PubKey:           validator.PubKey,
		ProposerPriority: validator.ProposerPriority,
		VotingPower:      validator.VotingPower,
	}, nil
}

// Validator Set at a height REST handler
func ValidatorSetRequestHandlerFn(clientCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse pagination parameters")
			return
		}

		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse block height")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse chain height")
			return
		}
		if height > chainHeight {
			rest.WriteErrorRes(w, "requested block height is bigger then the chain length")
			return
		}

		output, err := GetValidators(clientCtx, &height, &page, &limit)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}

