package helper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	tmclient "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
	"github.com/pkg/errors"
	"strconv"
)
var (
	errInitWrongFlags = errors.New("expected either (--hash/-x & --height) OR --force/-f, none given")
)
// InitLight is a helper function for init light
func InitLight(chain *collactor.Chain, force bool) (string, error) {
	db, df, err := chain.NewLightDB()
	if err != nil {
		return "", err
	}
	defer df()

	switch {
	case force: // force initialization from trusted node
		_, err := chain.LightClientWithoutTrust(db)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("successfully created light client for %s by trusting endpoint %s...",
			chain.ChainID, chain.RPCAddr), nil
	//case height > 0 && len(hash) > 0: // height and hash are given
	//	_, err = chain.LightClientWithTrust(db, chain.TrustOptions(height, hash))
	//	if err != nil {
	//		return "", wrapInitFailed(err)
	//	}
	//	return "", nil
	default: // return error
		return "", errInitWrongFlags
	}
}

// GetLightHeader returns header with chain and optional height as inputs
func GetLightHeader(chain *collactor.Chain, opts ...string) (*tmclient.Header, error) {
	if len(opts) > 0 {
		height, err := strconv.ParseInt(opts[0], 10, 64) //convert to int64
		if err != nil {
			return nil, err
		}

		if height <= 0 {
			height, err = chain.GetLatestLightHeight()
			if err != nil {
				return nil, err
			}

			if height < 0 {
				return nil, collactor.ErrLightNotInitialized
			}
		}

		return chain.GetLightSignedHeaderAtHeight(height)
	}

	return chain.GetLatestLightHeader()
}