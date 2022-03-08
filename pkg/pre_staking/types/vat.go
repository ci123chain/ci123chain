package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"time"
)

type StakingVaultOld struct {
	ID 			 uint64		  `json:"id"`
	StartTime    time.Time    `json:"start_time"`
	StorageTime  time.Duration    `json:"storage_time"`
	EndTime      time.Time    `json:"end_time"`
	Amount       sdk.Coin     `json:"amount"`
	Validator 	 sdk.AccAddress	`json:"validator"`
	Delegator 	 sdk.AccAddress	`json:"delegator"`
	TransLogs 	 []sdk.AccAddress `json:"trans_logs"`
	Processed 	 bool		   `json:"processed"`
}

type StakingVault struct {
	ID 			 uint64		  `json:"id"`
	StartTime    time.Time    `json:"start_time"`
	StorageTime  time.Duration    `json:"storage_time"`
	EndTime      time.Time    `json:"end_time"`
	Amount       sdk.Coin     `json:"amount"`
	Validator 	 sdk.AccAddress	`json:"validator"`
	Delegator 	 sdk.AccAddress	`json:"delegator"`
	TransLogs 	 []TransLog `json:"trans_logs"`
	Processed 	 bool		   `json:"processed"`
}

type TransLog struct {
	Src sdk.AccAddress	`json:"src"`
	Dst sdk.AccAddress	`json:"dst"`
	LogTime time.Time	`json:"log_time"`
}

func StakingVaultFrom(svo StakingVaultOld) StakingVault {

	return StakingVault{
		ID: svo.ID,
		StartTime: svo.StartTime,
		StorageTime: svo.StorageTime,
		EndTime: svo.EndTime,
		Amount: svo.Amount,
		Validator: svo.Validator,
		Delegator: svo.Delegator,
		Processed: svo.Processed,
	}
}

func NewTransLog(src, dst sdk.AccAddress, logTime time.Time) TransLog {
	return TransLog{
		Src: src,
		Dst: dst,
		LogTime: logTime,
	}
}

func NewStakingVault(id uint64,sta, et time.Time, st time.Duration, amount sdk.Coin, val, del sdk.AccAddress) StakingVault {
	return StakingVault{
		ID: 		id,
		StartTime:   sta,
		StorageTime: st,
		EndTime:     et,
		Amount:      amount,
		Validator: val,
		Delegator: del,
	}
}
//
//type VaultRecord struct {
//	LatestVaultID       *big.Int    `json:"latest_vault_id"`
//	Vaults              map[string]StakingVault    `json:"vaults"`
//}
//
//func NewVaultRecord(id *big.Int,vaults map[string]StakingVault) VaultRecord {
//	return VaultRecord{
//		LatestVaultID: id,
//		Vaults:        vaults,
//	}
//}

//
//func NewEmptyVaultRecord() VaultRecord {
//	return VaultRecord{
//		LatestVaultID: nil,
//		Vaults:        nil,
//	}
//}
//
//func (vr VaultRecord) IsEmpty() bool {
//	return vr.LatestVaultID == nil
//}
//
//func (vr *VaultRecord) AddVault(v Vault) {
//	var latestId = new(big.Int).SetUint64(1)
//	if vr.IsEmpty() {
//		vr.Vaults = make(map[string]Vault, 0)
//		vr.LatestVaultID = new(big.Int).SetUint64(0)
//		vr.LatestVaultID = latestId
//	}else {
//		latestId.Add(vr.LatestVaultID, new(big.Int).SetUint64(1))
//		vr.LatestVaultID = latestId
//	}
//	vr.Vaults[latestId.String()] = v
//}
//
//func (vr *VaultRecord) PopVaultAmountAndEndTime(id *big.Int) (sdk.Coin, time.Time, error) {
//	if vr.LatestVaultID.Cmp(id) == -1 {
//		return sdk.NewEmptyCoin(), time.Time{}, errors.New("this no id match")
//	}
//	v := vr.Vaults[id.String()]
//	if !v.Amount.IsPositive() {
//		return sdk.NewEmptyCoin(), time.Time{}, errors.New("no balance to un_delegate in this vault")
//	}
//	oldVault := vr.Vaults[id.String()]
//	//delete(vr.Vaults, id.String())
//	vr.Vaults[id.String()] = NewVault(oldVault.StartTime, oldVault.EndTime, oldVault.StorageTime, sdk.NewEmptyCoin())
//	return oldVault.Amount, oldVault.EndTime, nil
//}