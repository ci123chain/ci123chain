package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

type AccAddress struct {
	common.Address
}

func ToAccAddress(addr []byte) AccAddress {
	return AccAddress{
		Address: common.BytesToAddress(addr),
	}
}

func HexToAddress(addres string) AccAddress {
	return AccAddress{
		common.HexToAddress(addres),
	}
}

func (aa AccAddress)Equal(aa2 AccAddress) bool {
	return aa.Address == aa2.Address
}

func (aa AccAddress) Empty() bool {
	return aa.Address == common.Address{}
}

func (aa AccAddress) Marshal() ([]byte, error) {
	return aa.Address.Bytes(), nil
}

//func (aa AccAddress) Validate() (error) {
//	if len(aa.Address) != 0 {
//		return errors.New("cannot override BaseAccount address")
//	}
//	return nil
//}

func (aa AccAddress) String() string {
	return aa.Address.Hex()
}


func (aa AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

func (aa *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	addr2 := common.HexToAddress(s)
	*aa = AccAddress{
		addr2,
	}
	return nil
}