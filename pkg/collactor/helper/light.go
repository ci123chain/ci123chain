package helper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/pkg/errors"
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
