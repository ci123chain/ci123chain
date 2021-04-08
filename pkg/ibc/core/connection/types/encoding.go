package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

func UnmarshalConnectionRequest (cdc *codec.Codec, bz []byte) (QueryConnectionsRequest, error) {
	var req QueryConnectionsRequest
	err := cdc.UnmarshalJSON(bz, &req)
	return req, err
}

func MustMarshalConnectionResponse (cdc *codec.Codec, resp QueryConnectionsResponse) []byte {
	bz := cdc.MustMarshalJSON(resp)
	return bz
}

func MustUnmarshalConnectionResponse (cdc *codec.Codec, bz []byte) QueryConnectionsResponse {
	var req QueryConnectionsResponse
	cdc.MustUnmarshalJSON(bz, &req)
	return req
}

