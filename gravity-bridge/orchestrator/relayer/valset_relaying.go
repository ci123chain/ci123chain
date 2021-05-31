package relayer

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

func relayValsets(
	currentValSet types.ValSet,
	ethKey *ecdsa.PrivateKey,
	client *jsonrpc.Client,
	contact cosmos_gravity.Contact,
	contractAddr, gravityId string,
	timeout time.Duration,
) {
	lg := logger.GetLogger()
	getLatestValsets := gravity_utils.Exec(func() interface{} {
		latestValsets, err := cosmos_gravity.GetLatestValSets(contact)
		if err != nil {
			return err
		}
		return latestValsets
	}).Await()

	latestValsets, ok := getLatestValsets.([]*types.ValSet)
	if !ok {
		lg.Info("Failed to get latest valsets")
		return
	}

	if latestValsets == nil {
		return
	}

	var latestConfirmed []types.ValsetConfirmResponse
	var latestValset types.ValSet
	var lastError error
	latestNonce := latestValsets[0].Nonce
	for {
		if latestNonce == 0 {
			break
		}
		getValSet := gravity_utils.Exec(func() interface{} {
			valset, err := cosmos_gravity.GetValSet(contact, latestNonce)
			if err != nil {
				return err
			}
			return valset
		}).Await()
		valset, ok := getValSet.(*types.ValSet)
		if ok {
			//assert_eq
			if latestNonce != valset.Nonce {
				panic("latestNonce not equal valsetNonce")
			}
			getConfirms := gravity_utils.Exec(func() interface{} {
				confirms, err := cosmos_gravity.GetAllValsetConfirms(contact, latestNonce)
				if err != nil {
					return err
				}
				return confirms
			}).Await()
			confirms, ok := getConfirms.([]types.ValsetConfirmResponse)
			if ok {
				var sigs []types.Confirm
				for _, confirm := range confirms {
					//assert_eq
					if valset.Nonce != confirm.Nonce {
						panic("confirmNonce not equal valsetNonce")
					}
					sigs = append(sigs, confirm)
				}
				hash := types.EncodeValsetConfirmHashed(gravityId, *valset)
				_, err := currentValSet.OrderSigs(hash, sigs)
				if err != nil {
					lastError = err
				}
				latestConfirmed = confirms
				latestValset = *valset
			}
		}

		latestNonce -= 1
	}

	if len(latestConfirmed) == 0 {
		lg.Error("We don't have a latest confirmed valset?")
		return;
	}

	latestCosmosValset := latestValset
	latestCosmosConfirmed := latestConfirmed

	if latestNonce > latestCosmosValset.Nonce && lastError != nil {
		lg.Error(lastError.Error())
	}

	latestCosmosValsetNonce := latestCosmosValset.Nonce
	if latestCosmosValsetNonce > currentValSet.Nonce {
		cost, err := ethereum_gravity.EstimateValsetCost(latestCosmosValset, currentValSet, latestCosmosConfirmed, client, contractAddr, gravityId, ethKey)
		if err != nil {
			lg.Error(fmt.Sprintf("Valset cost estimate for Nonce %d failed, err: %s",latestCosmosValset.Nonce, err.Error()))
			return
		}

		lg.Info(fmt.Sprintf("We have detected latest valset %d but latest on Ethereum is %d This valset will be submit, expected cost: %v", latestCosmosValset.Nonce, currentValSet.Nonce, cost))
		gravity_utils.Exec(func() interface{} {
			if err := ethereum_gravity.SendEthValsetUpdate(latestCosmosValset, currentValSet, latestCosmosConfirmed, client, timeout, contractAddr, gravityId, ethKey); err != nil {
				return err
			}
			return nil
		}).Await()
	}
}