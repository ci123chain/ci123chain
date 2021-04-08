package collactor

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	clientutils "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/utils"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	ibctmtypes "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/light"
	"reflect"
	"time"
)

// CreateClients creates clients for src on dst and dst on src if the client ids are unspecified.
// TODO: de-duplicate code
func (c *Chain) CreateClients(dst *Chain) (modified bool, err error) {
	// Handle off chain light clients
	if err := c.ValidateLightInitialized(); err != nil {
		return false, err
	}

	if err = dst.ValidateLightInitialized(); err != nil {
		return false, err
	}

	srcUpdateHeader, dstUpdateHeader, err := GetIBCCreateClientHeaders(c, dst)
	if err != nil {
		return false, err
	}

	// Create client for the destination chain on the source chain if client id is unspecified
	if c.PathEnd.ClientID == "" {
		if c.debug {
			c.logCreateClient(dst, dstUpdateHeader.Header.Height)
		}
		ubdPeriod, err := dst.QueryUnbondingPeriod()
		if err != nil {
			return modified, err
		}

		// Create the ClientState we want on 'c' tracking 'dst'
		clientState := ibctmtypes.NewClientState(
			dstUpdateHeader.GetHeader().GetChainID(),
			ibctmtypes.NewFractionFromTm(light.DefaultTrustLevel),
			dst.GetTrustingPeriod(),
			ubdPeriod,
			time.Minute*10,
			dstUpdateHeader.GetHeight().(clienttypes.Height),
			commitmenttypes.GetSDKSpecs(),
			DefaultUpgradePath,
			AllowUpdateAfterExpiry,
			AllowUpdateAfterMisbehaviour,
		)

		// Check if an identical light client already exists
		clientID, found := FindMatchingClient(c, dst, clientState)
		if !found {
			msgs := []sdk.Msg{
				c.CreateClient(
					clientState,
					dstUpdateHeader,
				),
			}

			// if a matching client does not exist, create one
			res, success, err := c.SendMsgs(msgs)
			if err != nil {
				return modified, err
			}
			if !success {
				return modified, fmt.Errorf("tx failed: %s", res.RawLog)
			}

			// update the client identifier
			// use index 0, the transaction only has one message
			clientID, err = ParseClientIDFromEvents(res.Logs[0].Events)
			if err != nil {
				return modified, err
			}
		}

		c.PathEnd.ClientID = clientID
		modified = true

	} else {
		// Ensure client exists in the event of user inputted identifiers
		// TODO: check client is not expired
		_, err := c.QueryClientState(srcUpdateHeader.Header.Height)
		if err != nil {
			return false, fmt.Errorf("please ensure provided on-chain client (%s) exists on the chain (%s): %v",
				c.PathEnd.ClientID, c.ChainID, err)
		}
	}

	// Create client for the source chain on destination chain if client id is unspecified
	if dst.PathEnd.ClientID == "" {
		if dst.debug {
			dst.logCreateClient(c, srcUpdateHeader.Header.Height)
		}
		ubdPeriod, err := c.QueryUnbondingPeriod()
		if err != nil {
			return modified, err
		}
		// Create the ClientState we want on 'dst' tracking 'c'
		clientState := ibctmtypes.NewClientState(
			srcUpdateHeader.GetHeader().GetChainID(),
			ibctmtypes.NewFractionFromTm(light.DefaultTrustLevel),
			c.GetTrustingPeriod(),
			ubdPeriod,
			time.Minute*10,
			srcUpdateHeader.GetHeight().(clienttypes.Height),
			commitmenttypes.GetSDKSpecs(),
			DefaultUpgradePath,
			AllowUpdateAfterExpiry,
			AllowUpdateAfterMisbehaviour,
		)

		// Check if an identical light client already exists
		// NOTE: we pass in 'dst' as the source and 'c' as the
		// counterparty.
		clientID, found := FindMatchingClient(dst, c, clientState)
		if !found {
			msgs := []sdk.Msg{
				dst.CreateClient(
					clientState,
					srcUpdateHeader,
				),
			}

			// if a matching client does not exist, create one
			res, success, err := dst.SendMsgs(msgs)
			if err != nil {
				return modified, err
			}
			if !success {
				return modified, fmt.Errorf("tx failed: %s", res.RawLog)
			}

			// update client identifier
			clientID, err = ParseClientIDFromEvents(res.Logs[0].Events)
			if err != nil {
				return modified, err
			}
		}
		dst.PathEnd.ClientID = clientID
		modified = true

	} else {
		// Ensure client exists in the event of user inputted identifiers
		// TODO: check client is not expired
		_, err := dst.QueryClientState(dstUpdateHeader.Header.Height)
		if err != nil {
			return false, fmt.Errorf("please ensure provided on-chain client (%s) exists on the chain (%s): %v",
				dst.PathEnd.ClientID, dst.ChainID, err)
		}

	}

	c.Log(fmt.Sprintf("★ Clients created: client(%s) on chain[%s] and client(%s) on chain[%s]",
		c.PathEnd.ClientID, c.ChainID, dst.PathEnd.ClientID, dst.ChainID))

	return modified, nil
}



// FindMatchingClient will determine if there exists a client with identical client and consensus states
// to the client which would have been created. Source is the chain that would be adding a client
// which would track the counterparty. Therefore we query source for the existing clients
// and check if any match the counterparty. The counterparty must have a matching consensus state
// to the latest consensus state of a potential match. The provided client state is the client
// state that will be created if there exist no matches.
func FindMatchingClient(source, counterparty *Chain, clientState *ibctmtypes.ClientState) (string, bool) {
	// TODO: add appropriate offset and limits, along with retries
	clientsResp, err := source.QueryClients(0, 1000)
	if err != nil {
		if source.debug {
			source.Log(fmt.Sprintf("Error: querying clients on %s failed: %v", source.PathEnd.ChainID, err))
		}
		return "", false
	}

	for _, identifiedClientState := range clientsResp.ClientStates {
		// unpack any into ibc tendermint client state
		clientStateExported := identifiedClientState.ClientState

		// cast from interface to concrete type
		existingClientState, ok := clientStateExported.(*ibctmtypes.ClientState)
		if !ok {
			return "", false
		}

		// check if the client states match
		// NOTE: IsFrozen is a sanity check, the client to be created should always
		// have a zero frozen height and therefore should never match with a frozen client
		if IsMatchingClient(*clientState, *existingClientState) && !existingClientState.IsFrozen() {

			// query the latest consensus state of the potential matching client
			consensusStateResp, err := clientutils.QueryConsensusStateABCI(source.CLIContext(0),
				identifiedClientState.ClientId, existingClientState.GetLatestHeight())
			if err != nil {
				if source.debug {
					source.Log(fmt.Sprintf("Error: failed to query latest consensus state for existing client on chain %s: %v",
						source.PathEnd.ChainID, err))
				}
				continue
			}

			//nolint:lll
			header, err := counterparty.GetLightSignedHeaderAtHeight(int64(existingClientState.GetLatestHeight().GetRevisionHeight()))
			if err != nil {
				if source.debug {
					source.Log(fmt.Sprintf("Error: failed to query header for chain %s at height %d: %v",
						counterparty.PathEnd.ChainID, existingClientState.GetLatestHeight().GetRevisionHeight(), err))
				}
				continue
			}

			exportedConsState := consensusStateResp.ConsensusState
			existingConsensusState, ok := exportedConsState.(*ibctmtypes.ConsensusState)
			if !ok {
				if source.debug {
					source.Log(fmt.Sprintf("Error:consensus state is not tendermint type on chain %s", counterparty.PathEnd.ChainID))
				}
				continue
			}

			if existingClientState.IsExpired(existingConsensusState.Timestamp, time.Now()) {
				continue
			}

			if IsMatchingConsensusState(existingConsensusState, header.ConsensusState()) {
				// found matching client
				return identifiedClientState.ClientId, true
			}
		}
	}

	return "", false
}

// IsMatchingClient determines if the two provided clients match in all fields
// except latest height. They are assumed to be IBC tendermint light clients.
// NOTE: we don't pass in a pointer so upstream references don't have a modified
// latest height set to zero.
func IsMatchingClient(clientStateA, clientStateB ibctmtypes.ClientState) bool {
	// zero out latest client height since this is determined and incremented
	// by on-chain updates. Changing the latest height does not fundamentally
	// change the client. The associated consensus state at the latest height
	// determines this last check
	clientStateA.LatestHeight = clienttypes.ZeroHeight()
	clientStateB.LatestHeight = clienttypes.ZeroHeight()

	return reflect.DeepEqual(clientStateA, clientStateB)
}

// IsMatchingConsensusState determines if the two provided consensus states are
// identical. They are assumed to be IBC tendermint light clients.
func IsMatchingConsensusState(consensusStateA, consensusStateB *ibctmtypes.ConsensusState) bool {
	return reflect.DeepEqual(*consensusStateA, *consensusStateB)
}