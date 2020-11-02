package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgStoreContent struct {
	FromAddress sdk.AccAddress  `json:"from_address"`
	Key         string          `json:"key"`
	Content     []byte   `json:"content"`
}


func NewMsgStoreContent(from sdk.AccAddress, key string, content []byte) MsgStoreContent {
	return MsgStoreContent{
		FromAddress: from,
		Key:  key,
		Content: content,
	}
}

func (msg MsgStoreContent) Route() string { return RouteKey }


func (msg MsgStoreContent) MsgType() string { return "store-content" }

func (msg MsgStoreContent) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgStoreContent) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg MsgStoreContent) Bytes() []byte {
	bytes, err := InfrastructureCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

type StoredContent struct {
	Key     string     `json:"key"`
	Content string     `json:"content"`
}

func NewStoredContent(key, content string) StoredContent{
	return StoredContent{
		Key:     key,
		Content: content,
	}
}