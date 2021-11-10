package types

import sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

type UpgradeHandler func(ctx sdk.Context, info []byte)

