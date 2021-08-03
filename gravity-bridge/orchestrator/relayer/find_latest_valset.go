package relayer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"sort"
)

const BLOCKS_TO_SEARCH = 5000

// This function finds the latest valset on the Gravity contract by looking back through the event
// history and finding the most recent ValsetUpdatedEvent. Most of the time this will be very fast
// as the latest update will be in recent blockchain history and the search moves from the present
// backwards in time. In the case that the validator set has not been updated for a very long time
// this will take longer.
func findLatestValset(contact cosmos_gravity.Contact, contractAddr string, client *jsonrpc.Client, ourEthereumAddress common.Address) (*types.ValSet, error) {
	getBlock := gravity_utils.Exec(func() interface{} {
		block, err := client.Eth().BlockNumber()
		if err != nil {
			return err
		}
		return block
	}).Await()

	currentBlock, ok := getBlock.(uint64)
	if !ok {
		return nil, getBlock.(error)
	}

	getLatestEthereumValSetNonce := gravity_utils.Exec(func() interface{} {
		nonce, err := ethereum_gravity.GetValSetNonce(contractAddr, ourEthereumAddress, client)
		if err != nil {
			return err
		}
		return nonce
	}).Await()

	latestEthereumValSetNonce, ok := getLatestEthereumValSetNonce.(uint64)
	if !ok {
		return nil, getLatestEthereumValSetNonce.(error)
	}

	getCosmosChainValSet := gravity_utils.Exec(func() interface{} {
		valSet, err := cosmos_gravity.GetValSet(contact, latestEthereumValSetNonce)
		if err != nil {
			return err
		}
		return valSet
	}).Await()

	cosmosChainValSet, ok := getCosmosChainValSet.(*types.ValSet)
	if !ok {
		return nil, getCosmosChainValSet.(error)
	}

	lg := logger.GetLogger()
	for {
		if currentBlock == 0 {
			break
		}
		lg.Info(fmt.Sprintf("About to submit a Valset or Batch looking back into the history to find the last Valset Update, on block %d", currentBlock))

		var endSearch uint64
		if currentBlock < BLOCKS_TO_SEARCH {
			endSearch = 0
		} else {
			endSearch = currentBlock - BLOCKS_TO_SEARCH
		}

		getAllValSetEvents := gravity_utils.Exec(func() interface{} {
			allValSetEvents, err := ethereum_gravity.CheckForEvents(endSearch, currentBlock, []string{contractAddr}, []string{"ValsetUpdatedEvent(uint256,address[],uint256[])"}, client)
			if err != nil {
				return err
			}
			return allValSetEvents
		}).Await()

		allValSetEvents, ok := getAllValSetEvents.([]*web3.Log)
		if !ok {
			return nil, getAllValSetEvents.(error)
		}

		// by default the lowest found valset goes first, we want the highest.
		for i, j := 0, len(allValSetEvents)-1; i < j; i, j = i+1, j-1 {
			allValSetEvents[i], allValSetEvents[j] = allValSetEvents[j], allValSetEvents[i]
		}

		lg.Info(fmt.Sprintf("Found events"))

		ethValSet := new(types.ValSet)
		if len(allValSetEvents) != 0 {
			valSetUpdatedEvent, err := types.ValSetUpdatedEventFromLog(allValSetEvents[0])
			if err != nil {
				return nil, err
			}
			ethValSet.Nonce = valSetUpdatedEvent.Nonce
			ethValSet.Members = valSetUpdatedEvent.Members
			same := checkIfValsetsDiffer(cosmosChainValSet, ethValSet)
			if !same {
				lg.Info(fmt.Sprintf("Validator sets for nonce: %d, Cosmos and Ethereum differ. Possible bridge highjacking!", valSetUpdatedEvent.Nonce))
			}
			return ethValSet, nil
		}
		currentBlock = endSearch;
	}

	lg.Error("Could not find the last validator set for contract %s, probably not a valid Gravity contract!", contractAddr)

	return nil, errors.New("Could not find the last validator set")
}

// This function exists to provide a warning if Cosmos and Ethereum have different validator sets
// for a given nonce. In the mundane version of this warning the validator sets disagree on sorting order
// which can happen if some relayer uses an unstable sort, or in a case of a mild griefing attack.
// The Gravity contract validates signatures in order of highest to lowest power. That way it can exit
// the loop early once a vote has enough power, if a relayer where to submit things in the reverse order
// they could grief users of the contract into paying more in gas.
// The other (and far worse) way a disagreement here could occur is if validators are colluding to steal
// funds from the Gravity contract and have submitted a highjacking update. If slashing for off Cosmos chain
// Ethereum signatures is implemented you would put that handler here.
func checkIfValsetsDiffer(cosmosValset, ethereumValset *types.ValSet) bool {
	lg := logger.GetLogger()
	if cosmosValset == nil && ethereumValset.Nonce == 0 {
		// bootstrapping case
		return true
	} else if cosmosValset == nil {
		lg.Error(fmt.Sprintf("Cosmos does not have a valset for nonce: %d, but that is the one on the Ethereum chain! Possible bridge highjacking!", ethereumValset.Nonce))
		return false
	}

	//?
	if cosmosValset.Nonce != ethereumValset.Nonce {
		lg.Error(fmt.Sprintf("Cosmos has the wrong validator set for nonce: %d. Possible bridge highjacking!", ethereumValset.Nonce))
		return false
	}

	cValSet := cosmosValset.Members
	eValSet := ethereumValset.Members

	sort.SliceStable(cValSet, func(i, j int) bool {
		if cValSet[i].Power > cValSet[j].Power {
			return true
		} else if cValSet[i].Power == cValSet[j].Power {
			if bytes.Compare(cValSet[i].EthAddress.Bytes(), cValSet[j].EthAddress.Bytes()) > 0 {
				return true
			}
		}
		return false
	})

	sort.SliceStable(eValSet, func(i, j int) bool {
		if eValSet[i].Power > eValSet[j].Power {
			return true
		} else if eValSet[i].Power == eValSet[j].Power {
			if bytes.Compare(eValSet[i].EthAddress.Bytes(), eValSet[j].EthAddress.Bytes()) > 0 {
				return true
			}
		}
		return false
	})

	//compare
	if len(cValSet) != len(eValSet) {
		lg.Error(fmt.Sprintf("Validator sets for nonce: %d, Cosmos and Ethereum differ. Possible bridge highjacking!", ethereumValset.Nonce))
		return false
	}

	for i := 0; i < len(cValSet); i++ {
		if cValSet[i].EthAddress.String() != eValSet[i].EthAddress.String() || cValSet[i].Power != eValSet[i].Power {
			return false
		}
	}

	return true
}


