// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: weelink/abci/result.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/tendermint/tendermint/abci/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Result struct {
	Code      uint32        `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Codespace string        `protobuf:"bytes,2,opt,name=codespace,proto3" json:"codespace,omitempty"`
	Data      []byte        `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Log       string        `protobuf:"bytes,4,opt,name=log,proto3" json:"log,omitempty"`
	GasWanted uint64        `protobuf:"varint,5,opt,name=gasWanted,proto3" json:"gasWanted,omitempty"`
	GasUsed   uint64        `protobuf:"varint,6,opt,name=gasUsed,proto3" json:"gasUsed,omitempty"`
	Events    []types.Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_1af8ab3cc216d1f1, []int{0}
}
func (m *Result) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Result.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return m.Size()
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Result) GetCodespace() string {
	if m != nil {
		return m.Codespace
	}
	return ""
}

func (m *Result) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Result) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *Result) GetGasWanted() uint64 {
	if m != nil {
		return m.GasWanted
	}
	return 0
}

func (m *Result) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *Result) GetEvents() []types.Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterType((*Result)(nil), "weelink.base.abci.Result")
}

func init() { proto.RegisterFile("weelink/abci/result.proto", fileDescriptor_1af8ab3cc216d1f1) }

var fileDescriptor_1af8ab3cc216d1f1 = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0x2b, 0x31,
	0x14, 0x86, 0x27, 0xb7, 0x73, 0xa7, 0x34, 0x2a, 0x68, 0x10, 0x89, 0x55, 0xc6, 0xc1, 0xd5, 0xac,
	0x32, 0xd8, 0xfa, 0x04, 0x05, 0x17, 0x6e, 0x03, 0x22, 0xb8, 0xcb, 0x64, 0x0e, 0x69, 0x68, 0x9b,
	0x0c, 0x4d, 0xaa, 0xf8, 0x16, 0x3e, 0x56, 0x57, 0xd2, 0xa5, 0x2b, 0x91, 0xf6, 0x45, 0x24, 0x69,
	0xa5, 0x5d, 0xcd, 0x37, 0xe7, 0xfb, 0xf3, 0xc3, 0x39, 0xf8, 0xf2, 0x0d, 0x60, 0xaa, 0xcd, 0xa4,
	0x12, 0xb5, 0xd4, 0xd5, 0x1c, 0xdc, 0x62, 0xea, 0x59, 0x3b, 0xb7, 0xde, 0x92, 0xb3, 0x9d, 0x62,
	0xb5, 0x70, 0xc0, 0x82, 0xef, 0x9f, 0x2b, 0xab, 0x6c, 0xb4, 0x55, 0xa0, 0x6d, 0xb0, 0x7f, 0xe5,
	0xc1, 0x34, 0x30, 0x9f, 0x69, 0xe3, 0xb7, 0x35, 0xfe, 0xbd, 0x05, 0xb7, 0x95, 0xb7, 0x9f, 0x08,
	0x67, 0x3c, 0xd6, 0x12, 0x82, 0x53, 0x69, 0x1b, 0xa0, 0xa8, 0x40, 0xe5, 0x09, 0x8f, 0x4c, 0xae,
	0x71, 0x2f, 0x7c, 0x5d, 0x2b, 0x24, 0xd0, 0x7f, 0x05, 0x2a, 0x7b, 0x7c, 0x3f, 0x08, 0x2f, 0x1a,
	0xe1, 0x05, 0xed, 0x14, 0xa8, 0x3c, 0xe6, 0x91, 0xc9, 0x29, 0xee, 0x4c, 0xad, 0xa2, 0x69, 0xcc,
	0x06, 0x0c, 0x1d, 0x4a, 0xb8, 0x67, 0x61, 0x3c, 0x34, 0xf4, 0x7f, 0x81, 0xca, 0x94, 0xef, 0x07,
	0x84, 0xe2, 0xae, 0x12, 0xee, 0xc9, 0x41, 0x43, 0xb3, 0xe8, 0xfe, 0x7e, 0xc9, 0x3d, 0xce, 0xe0,
	0x15, 0x8c, 0x77, 0xb4, 0x5b, 0x74, 0xca, 0xa3, 0xc1, 0x05, 0xdb, 0x2f, 0x12, 0xf7, 0x65, 0x0f,
	0x41, 0x8f, 0xd2, 0xe5, 0xf7, 0x4d, 0xc2, 0x77, 0xd9, 0xd1, 0xe3, 0x72, 0x9d, 0xa3, 0xd5, 0x3a,
	0x47, 0x3f, 0xeb, 0x1c, 0x7d, 0x6c, 0xf2, 0x64, 0xb5, 0xc9, 0x93, 0xaf, 0x4d, 0x9e, 0xbc, 0x54,
	0x4a, 0xfb, 0xf1, 0xa2, 0x66, 0xd2, 0xce, 0x2a, 0xa9, 0xef, 0x06, 0x43, 0x39, 0x16, 0xda, 0x1c,
	0x62, 0x3b, 0x51, 0x07, 0x17, 0xaa, 0xb3, 0x78, 0xa2, 0xe1, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xb3, 0x20, 0x10, 0xab, 0x85, 0x01, 0x00, 0x00,
}

func (m *Result) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Result) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Result) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintResult(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.GasUsed != 0 {
		i = encodeVarintResult(dAtA, i, uint64(m.GasUsed))
		i--
		dAtA[i] = 0x30
	}
	if m.GasWanted != 0 {
		i = encodeVarintResult(dAtA, i, uint64(m.GasWanted))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Log) > 0 {
		i -= len(m.Log)
		copy(dAtA[i:], m.Log)
		i = encodeVarintResult(dAtA, i, uint64(len(m.Log)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintResult(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Codespace) > 0 {
		i -= len(m.Codespace)
		copy(dAtA[i:], m.Codespace)
		i = encodeVarintResult(dAtA, i, uint64(len(m.Codespace)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintResult(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintResult(dAtA []byte, offset int, v uint64) int {
	offset -= sovResult(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Result) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovResult(uint64(m.Code))
	}
	l = len(m.Codespace)
	if l > 0 {
		n += 1 + l + sovResult(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovResult(uint64(l))
	}
	l = len(m.Log)
	if l > 0 {
		n += 1 + l + sovResult(uint64(l))
	}
	if m.GasWanted != 0 {
		n += 1 + sovResult(uint64(m.GasWanted))
	}
	if m.GasUsed != 0 {
		n += 1 + sovResult(uint64(m.GasUsed))
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovResult(uint64(l))
		}
	}
	return n
}

func sovResult(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozResult(x uint64) (n int) {
	return sovResult(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Result) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowResult
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Result: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Result: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Codespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResult
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Codespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthResult
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Log", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResult
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Log = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasWanted", wireType)
			}
			m.GasWanted = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasWanted |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasUsed", wireType)
			}
			m.GasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthResult
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, types.Event{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipResult(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthResult
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipResult(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowResult
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowResult
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowResult
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthResult
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupResult
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthResult
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthResult        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowResult          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupResult = fmt.Errorf("proto: unexpected end of group")
)
