package collactor

import (
	"fmt"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	conntypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
)

func (c *Chain) logCreateClient(dst *Chain, dstH int64) {
	c.Log(fmt.Sprintf("- [%s] -> creating client on %s for %s header-height{%d} trust-period(%s)",
		c.ChainID, c.ChainID, dst.ChainID, dstH, dst.GetTrustingPeriod()))
}


func logConnectionStates(src, dst *Chain, srcConn, dstConn *conntypes.QueryConnectionResponse) {
	src.Log(fmt.Sprintf("- [%s]@{%d}conn(%s)-{%s} : [%s]@{%d}conn(%s)-{%s}",
		src.ChainID,
		MustGetHeight(srcConn.ProofHeight),
		src.PathEnd.ConnectionID,
		srcConn.Connection.State,
		dst.ChainID,
		MustGetHeight(dstConn.ProofHeight),
		dst.PathEnd.ConnectionID,
		dstConn.Connection.State,
	))
}


func logChannelStates(src, dst *Chain, srcChan, dstChan *chantypes.QueryChannelResponse) {
	src.Log(fmt.Sprintf("- [%s]@{%d}chan(%s)-{%s} : [%s]@{%d}chan(%s)-{%s}",
		src.ChainID,
		MustGetHeight(srcChan.ProofHeight),
		src.PathEnd.ChannelID,
		srcChan.Channel.State,
		dst.ChainID,
		MustGetHeight(dstChan.ProofHeight),
		dst.PathEnd.ChannelID,
		dstChan.Channel.State,
	))
}
