package types

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"math/big"
	"time"
)

type Vault struct {
	StartTime    time.Time    `json:"start_time"`
	StorageTime  time.Duration    `json:"storage_time"`
	EndTime      time.Time    `json:"end_time"`
	Amount       sdk.Coin     `json:"amount"`
}

func NewVault(sta, et time.Time, st time.Duration, amount sdk.Coin) Vault {
	return Vault{
		StartTime:   sta,
		StorageTime: st,
		EndTime:     et,
		Amount:      amount,
	}
}

type VaultRecord struct {
	LatestVaultID       *big.Int    `json:"latest_vault_id"`
	Vaults              map[string]Vault    `json:"vaults"`
}

func NewVaultRecord(id *big.Int,vaults map[string]Vault) VaultRecord {
	return VaultRecord{
		LatestVaultID: id,
		Vaults:        vaults,
	}
}


func NewEmptyVaultRecord() VaultRecord {
	return VaultRecord{
		LatestVaultID: nil,
		Vaults:        nil,
	}
}

func (vr VaultRecord) IsEmpty() bool {
	return vr.LatestVaultID == nil
}

func (vr *VaultRecord) AddVault(v Vault) {
	var latestId = new(big.Int).SetUint64(1)
	if vr.IsEmpty() {
		vr.Vaults = make(map[string]Vault, 0)
		vr.LatestVaultID = new(big.Int).SetUint64(0)
		vr.LatestVaultID = latestId
	}else {
		latestId.Add(vr.LatestVaultID, new(big.Int).SetUint64(1))
		vr.LatestVaultID = latestId
	}
	vr.Vaults[latestId.String()] = v
}

func (vr *VaultRecord) PopVaultAmountAndEndTime(id *big.Int) (sdk.Coin, time.Time, error) {
	if vr.LatestVaultID.Cmp(id) == -1 {
		return sdk.NewEmptyCoin(), time.Time{}, errors.New("this no id match")
	}
	v := vr.Vaults[id.String()]
	if !v.Amount.IsPositive() {
		return sdk.NewEmptyCoin(), time.Time{}, errors.New("no balance to un_delegate in this vault")
	}
	oldVault := vr.Vaults[id.String()]
	//delete(vr.Vaults, id.String())
	vr.Vaults[id.String()] = NewVault(oldVault.StartTime, oldVault.EndTime, oldVault.StorageTime, sdk.NewEmptyCoin())
	return oldVault.Amount, oldVault.EndTime, nil
}