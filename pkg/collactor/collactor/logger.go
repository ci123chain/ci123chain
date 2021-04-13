package collactor

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	conntypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"strings"
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


// LogFailedTx takes the transaction and the messages to create it and logs the appropriate data
func (c *Chain) LogFailedTx(res *sdk.TxResponse, err error, msgs []sdk.Msg) {
	if c.debug {
		c.Log(fmt.Sprintf("- [%s] -> sending transaction:", c.ChainID))
		for _, msg := range msgs {
			c.Print(msg, false, false)
		}
	}

	if err != nil {
		c.logger.Error(fmt.Errorf("- [%s] -> err(%v)", c.ChainID, err).Error())
		if res == nil {
			return
		}
	}

	if res.Code != 0 && res.Codespace != "" {
		c.logger.Info(fmt.Sprintf("✘ [%s]@{%d} - msg(%s) err(%s:%d:%s)",
			c.ChainID, res.Height, getMsgAction(msgs), res.Codespace, res.Code, res.RawLog))
	}

	if c.debug && !res.Empty() {
		c.Log("- transaction response:")
		c.Print(res, false, false)
	}
}

// LogSuccessTx take the transaction and the messages to create it and logs the appropriate data
func (c *Chain) LogSuccessTx(res *sdk.TxResponse, msgs []sdk.Msg) {
	c.logger.Info(fmt.Sprintf("✔ [%s]@{%d} - msg(%s) hash(%s)", c.ChainID, res.Height, getMsgAction(msgs), res.TxHash))
}

func getMsgAction(msgs []sdk.Msg) string {
	var out string
	for i, msg := range msgs {
		out += fmt.Sprintf("%d:%s,", i, msg.MsgType())
	}
	return strings.TrimSuffix(out, ",")
}

// Print fmt.Printlns the json or yaml representation of whatever is passed in
// CONTRACT: The cmd calling this function needs to have the "json" and "indent" flags set
// TODO: better "text" printing here would be a nice to have
// TODO: fix indenting all over the code base
func (c *Chain) Print(toPrint interface{}, text, indent bool) error {
	var (
		out []byte
		err error
	)

	switch {
	case indent && text:
		return fmt.Errorf("must pass either indent or text, not both")
	case text:
		// TODO: This isn't really a good option,
		out = []byte(fmt.Sprintf("%v", toPrint))
	default:
		out, err = c.cdc.MarshalJSON(toPrint)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}