package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
)

type QueryParamsRequest struct {
}

type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

// QuerySigningInfoRequest is the request type for the Query/SigningInfo RPC
// method
type QuerySigningInfoRequest struct {
	// cons_address is the address to query signing info of
	ConsAddress string `protobuf:"bytes,1,opt,name=cons_address,json=consAddress,proto3" json:"cons_address,omitempty"`
}

func (m *QuerySigningInfoRequest) GetConsAddress() string {
	if m != nil {
		return m.ConsAddress
	}
	return ""
}

type QuerySigningInfoResponse struct {
	// val_signing_info is the signing info of requested val cons address
	ValSigningInfo ValidatorSigningInfo `protobuf:"bytes,1,opt,name=val_signing_info,json=valSigningInfo,proto3" json:"val_signing_info"`
}

func (m *QuerySigningInfoResponse) GetValSigningInfo() ValidatorSigningInfo {
	if m != nil {
		return m.ValSigningInfo
	}
	return ValidatorSigningInfo{}
}

// QuerySigningInfosRequest is the request type for the Query/SigningInfos RPC
// method
type QuerySigningInfosRequest struct {
	Pagination *pagination.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QuerySigningInfosRequest) GetPagination() *pagination.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QuerySigningInfosResponse struct {
	// info is the signing info of all validators
	Info       []ValidatorSigningInfo `protobuf:"bytes,1,rep,name=info,proto3" json:"info"`
	Pagination *pagination.PageResponse    `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QuerySigningInfosResponse) GetInfo() []ValidatorSigningInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *QuerySigningInfosResponse) GetPagination() *pagination.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}