package types

// GenesisState defines the ibc-transfer genesis state
type GenesisState struct {
	PortId      string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
	DenomTraces Traces `protobuf:"bytes,2,rep,name=denom_traces,json=denomTraces,proto3,castrepeated=Traces" json:"denom_traces" yaml:"denom_traces"`
	Params      Params `protobuf:"bytes,3,opt,name=params,proto3" json:"params"`
}