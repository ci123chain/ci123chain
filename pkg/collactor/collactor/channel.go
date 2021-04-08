package collactor

import (
	"fmt"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	"time"
)

// CreateOpenChannels runs the channel creation messages on timeout until they pass
func (c *Chain) CreateOpenChannels(dst *Chain, maxRetries uint64, to time.Duration) (modified bool, err error) {
	// client and connection identifiers must be filled in
	if err := ValidateConnectionPaths(c, dst); err != nil {
		return modified, err
	}
	// ports must be valid and channel ORDER must be the same
	if err := ValidateChannelParams(c, dst); err != nil {
		return modified, err
	}

	ticker := time.NewTicker(to)
	failures := uint64(0)
	for ; true; <-ticker.C {
		success, lastStep, recentlyModified, err := ExecuteChannelStep(c, dst)
		if err != nil {
			c.Log(err.Error())
		}
		if recentlyModified {
			modified = true
		}

		switch {
		// In the case of success and this being the last transaction
		// debug logging, log created channel and break
		case success && lastStep:

			if c.debug {
				srch, dsth, err := GetLatestLightHeights(c, dst)
				if err != nil {
					return modified, err
				}
				srcChan, dstChan, err := QueryChannelPair(c, dst, srch, dsth)
				if err != nil {
					return modified, err
				}
				logChannelStates(c, dst, srcChan, dstChan)
			}

			c.Log(fmt.Sprintf("★ Channel created: [%s]chan{%s}port{%s} -> [%s]chan{%s}port{%s}",
				c.ChainID, c.PathEnd.ChannelID, c.PathEnd.PortID,
				dst.ChainID, dst.PathEnd.ChannelID, dst.PathEnd.PortID))
			return modified, nil

		// In the case of success, reset the failures counter
		case success:
			failures = 0
			continue

		// In the case of failure, increment the failures counter and exit if this is the 3rd failure
		case !success:
			failures++
			c.Log("retrying transaction...")
			time.Sleep(5 * time.Second)

			if failures > maxRetries {
				return modified, fmt.Errorf("! Channel failed: [%s]chan{%s}port{%s} -> [%s]chan{%s}port{%s}",
					c.ChainID, c.PathEnd.ChannelID, c.PathEnd.PortID,
					dst.ChainID, dst.PathEnd.ChannelID, dst.PathEnd.PortID)
			}
		}
	}

	return modified, nil // lgtm [go/unreachable-statement]
}


// ExecuteChannelStep executes the next channel step based on the
// states of two channel ends specified by the relayer configuration
// file. The booleans return indicate if the message was successfully
// executed and if this was the last handshake step.
func ExecuteChannelStep(src, dst *Chain) (success, last, modified bool, err error) {
	if _, _, err := UpdateLightClients(src, dst); err != nil {
		return false, false, false, err
	}

	// if either identifier is missing, an existing channel that matches the required fields
	// is chosen or a new channel is created.
	if src.PathEnd.ChannelID == "" || dst.PathEnd.ChannelID == "" {
		success, modified, err := InitializeChannel(src, dst)
		if err != nil {
			return false, false, false, err
		}

		return success, false, modified, nil
	}

	// Query Channel data from src and dst
	srcChan, dstChan, err := QueryChannelPair(src, dst, int64(src.MustGetLatestLightHeight())-1,
		int64(dst.MustGetLatestLightHeight()-1))
	if err != nil {
		return false, false, false, err
	}

	switch {

	// OpenTry on source in case of crossing hellos (both channels are on INIT)
	// obtain proof of counterparty in TRYOPEN state and submit to source chain to update state
	// from INIT to TRYOPEN.
	case srcChan.Channel.State == chantypes.INIT && dstChan.Channel.State == chantypes.INIT:
		if src.debug {
			logChannelStates(src, dst, srcChan, dstChan)
		}

		msgs, err := src.ChanTry(dst)
		if err != nil {
			return false, false, false, err
		}

		_, success, err = src.SendMsgs(msgs)
		if !success {
			return false, false, false, err
		}

	// OpenAck on source if dst is at TRYOPEN and src is at INIT or TRYOPEN (crossing hellos)
	// obtain proof of counterparty in TRYOPEN state and submit to source chain to update state
	// from INIT/TRYOPEN to OPEN.
	case (srcChan.Channel.State == chantypes.INIT ||
		srcChan.Channel.State == chantypes.TRYOPEN) && dstChan.Channel.State == chantypes.TRYOPEN:
		if src.debug {
			logChannelStates(src, dst, srcChan, dstChan)
		}

		msgs, err := src.ChanAck(dst)
		if err != nil {
			return false, false, false, err
		}

		_, success, err = src.SendMsgs(msgs)
		if !success {
			return false, false, false, err
		}

	// OpenAck on counterparty
	// obtain proof of source in TRYOPEN state and submit to counterparty chain to update state
	// from INIT to OPEN.
	case srcChan.Channel.State == chantypes.TRYOPEN && dstChan.Channel.State == chantypes.INIT:
		if dst.debug {
			logChannelStates(dst, src, dstChan, srcChan)
		}

		msgs, err := dst.ChanAck(src)
		if err != nil {
			return false, false, false, err
		}

		_, success, err = dst.SendMsgs(msgs)
		if !success {
			return false, false, false, err
		}

	// OpenConfirm on source
	case srcChan.Channel.State == chantypes.TRYOPEN && dstChan.Channel.State == chantypes.OPEN:
		if src.debug {
			logChannelStates(src, dst, srcChan, dstChan)
		}

		msgs, err := src.ChanConfirm(dst)
		if err != nil {
			return false, false, false, err
		}

		last = true

		_, success, err = src.SendMsgs(msgs)
		if !success {
			return false, false, false, err
		}

	// OpenConfrim on counterparty
	case srcChan.Channel.State == chantypes.OPEN && dstChan.Channel.State == chantypes.TRYOPEN:
		if dst.debug {
			logChannelStates(dst, src, dstChan, srcChan)
		}

		msgs, err := dst.ChanConfirm(src)
		if err != nil {
			return false, false, false, err
		}

		_, success, err = dst.SendMsgs(msgs)
		if !success {
			return false, false, false, err
		}

		last = true

	case srcChan.Channel.State == chantypes.OPEN && dstChan.Channel.State == chantypes.OPEN:
		last = true

	}

	return true, last, false, nil
}

// InitializeChannel creates a new channel on either the source or destination chain .
// The identifiers set in the PathEnd's are used to determine which channel ends need to be
// initialized. The PathEnds are updated upon a successful transaction.
// NOTE: This function may need to be called twice if neither channel exists.
func InitializeChannel(src, dst *Chain) (success, modified bool, err error) {
	switch {

	// OpenInit on source
	// Neither channel has been initialized
	case src.PathEnd.ChannelID == "" && dst.PathEnd.ChannelID == "":
		//nolint:staticcheck
		if src.debug {
			// TODO: log that we are attempting to create new channel ends
		}

		channelID, found := FindMatchingChannel(src, dst)
		if !found {
			msgs, err := src.ChanInit(dst)
			if err != nil {
				return false, false, err
			}

			res, success, err := src.SendMsgs(msgs)
			if !success {
				return false, false, err
			}

			// update channel identifier in PathEnd
			// use index 1, channel open init is the second message in the transaction
			channelID, err = ParseChannelIDFromEvents(res.Logs[1].Events)
			if err != nil {
				return false, false, err
			}
		}
		src.PathEnd.ChannelID = channelID

		return true, true, nil

	// OpenTry on source
	// source channel does not exist, but counterparty channel exists
	case src.PathEnd.ChannelID == "" && dst.PathEnd.ChannelID != "":
		//nolint:staticcheck
		if src.debug {
			// TODO: update logging
		}

		channelID, found := FindMatchingChannel(src, dst)
		if !found {
			// open try on source chain
			msgs, err := src.ChanTry(dst)
			if err != nil {
				return false, false, err
			}

			res, success, err := src.SendMsgs(msgs)
			if !success {
				return false, false, err
			}

			// update channel identifier in PathEnd
			// use index 1, channel open try is the second message in the transaction
			channelID, err = ParseChannelIDFromEvents(res.Logs[1].Events)
			if err != nil {
				return false, false, err
			}
		}
		src.PathEnd.ChannelID = channelID

		return true, true, nil

	// OpenTry on counterparty
	// source channel exists, but counterparty channel does not exist
	case src.PathEnd.ChannelID != "" && dst.PathEnd.ChannelID == "":
		//nolint:staticcheck
		if dst.debug {
			// TODO: update logging
		}

		channelID, found := FindMatchingChannel(dst, src)
		if !found {
			// open try on destination chain
			msgs, err := dst.ChanTry(src)
			if err != nil {
				return false, false, err
			}

			res, success, err := dst.SendMsgs(msgs)
			if !success {
				return false, false, err
			}

			// update channel identifier in PathEnd
			// use index 1, channel open try is the second message in the transaction
			channelID, err = ParseChannelIDFromEvents(res.Logs[1].Events)
			if err != nil {
				return false, false, err
			}
		}
		dst.PathEnd.ChannelID = channelID

		return true, true, nil

	default:
		return false, false, fmt.Errorf("channel ends already created")
	}
}



// FindMatchingChannel will determine if there already exists a channel between source and counterparty
// that matches the parameters set in the relayer config.
func FindMatchingChannel(source, counterparty *Chain) (string, bool) {
	// TODO: add appropriate offset and limits, along with retries
	channelsResp, err := source.QueryChannels(0, 1000)
	if err != nil {
		if source.debug {
			source.Log(fmt.Sprintf("Error: querying channels on %s failed: %v", source.PathEnd.ChainID, err))
		}
		return "", false
	}

	for _, channel := range channelsResp.Channels {
		if IsMatchingChannel(source, counterparty, channel) {
			// unused channel found
			return channel.ChannelId, true
		}
	}

	return "", false
}

// IsMatchingChannel determines if given channel matches required conditions
func IsMatchingChannel(source, counterparty *Chain, channel *chantypes.IdentifiedChannel) bool {
	return channel.Ordering == source.PathEnd.GetOrder() &&
		IsConnectionFound(channel.ConnectionHops, source.PathEnd.ConnectionID) &&
		channel.Version == source.PathEnd.Version &&
		channel.PortId == source.PathEnd.PortID && channel.Counterparty.PortId == counterparty.PathEnd.PortID &&
		(((channel.State == chantypes.INIT || channel.State == chantypes.TRYOPEN) && channel.Counterparty.ChannelId == "") ||
			(channel.State == chantypes.OPEN && (counterparty.PathEnd.ChannelID == "" ||
				channel.Counterparty.ChannelId == counterparty.PathEnd.ChannelID)))
}