package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

func MustUnmarshalQueryChannelsResp(cdc *codec.Codec, bz []byte) QueryChannelsResponse {
	var resp QueryChannelsResponse
	cdc.MustUnmarshalJSON(bz, &resp)
	return resp
}

func MustMarshalQueryChannelsResp(cdc *codec.Codec, resp QueryChannelsResponse) []byte {
	return cdc.MustMarshalJSON(&resp)
}


func UnmarshalQueryChannelRequest(cdc *codec.Codec, bz []byte) (QueryChannelsRequest, error) {
	var resp QueryChannelsRequest
	err := cdc.UnmarshalJSON(bz, &resp)
	return resp, err
}