package types


// CapabilityOwners defines a set of owners of a single Capability. The set of
// owners must be unique.
type CapabilityOwners struct {
	Owners []Owner `protobuf:"bytes,1,rep,name=owners,proto3" json:"owners"`
}

type Owner struct {
	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty" yaml:"module"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" yaml:"name"`
}

type Capability struct {
	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty" yaml:"index"`
}


func (m *Capability) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}