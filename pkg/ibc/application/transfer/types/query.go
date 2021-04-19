package types

import "github.com/ci123chain/ci123chain/pkg/abci/types/pagination"

// QueryConnectionsRequest is the request type for the Query/DenomTraces RPC
// method
type QueryDenomTracesRequest struct {
	// pagination defines an optional pagination for the request.
	Pagination *pagination.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryConnectionsResponse is the response type for the Query/DenomTraces RPC
// method.
type QueryDenomTracesResponse struct {
	// denom_traces returns all denominations trace information.
	DenomTraces Traces `protobuf:"bytes,1,rep,name=denom_traces,json=denomTraces,proto3,castrepeated=Traces" json:"denom_traces"`
	// pagination defines the pagination in the response.
	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}