package types
//
//import (
//	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
//)
//
//type MsgChannelOpenInit struct {
//	PortId  string  `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
//	Channel Channel `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel"`
//	Signer  string  `protobuf:"bytes,3,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//
//
//// MsgChannelOpenInitResponse defines the Msg/ChannelOpenInit response type.
//type MsgChannelOpenInitResponse struct {
//}
//
//// MsgChannelOpenInit defines a msg sent by a Relayer to try to open a channel
//// on Chain B.
//type MsgChannelOpenTry struct {
//	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
//	// in the case of crossing hello's, when both chains call OpenInit, we need the channel identifier
//	// of the previous channel in state INIT
//	PreviousChannelId   string       `protobuf:"bytes,2,opt,name=previous_channel_id,json=previousChannelId,proto3" json:"previous_channel_id,omitempty" yaml:"previous_channel_id"`
//	Channel             Channel      `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel"`
//	CounterpartyVersion string       `protobuf:"bytes,4,opt,name=counterparty_version,json=counterpartyVersion,proto3" json:"counterparty_version,omitempty" yaml:"counterparty_version"`
//	ProofInit           []byte       `protobuf:"bytes,5,opt,name=proof_init,json=proofInit,proto3" json:"proof_init,omitempty" yaml:"proof_init"`
//	ProofHeight         clienttypes.Height `protobuf:"bytes,6,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	Signer              string       `protobuf:"bytes,7,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//// MsgChannelOpenTryResponse defines the Msg/ChannelOpenTry response type.
//type MsgChannelOpenTryResponse struct {
//}
//
//
//
//// MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
//// the change of channel state to TRYOPEN on Chain B.
//type MsgChannelOpenAck struct {
//	PortId                string       `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
//	ChannelId             string       `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty" yaml:"channel_id"`
//	CounterpartyChannelId string       `protobuf:"bytes,3,opt,name=counterparty_channel_id,json=counterpartyChannelId,proto3" json:"counterparty_channel_id,omitempty" yaml:"counterparty_channel_id"`
//	CounterpartyVersion   string       `protobuf:"bytes,4,opt,name=counterparty_version,json=counterpartyVersion,proto3" json:"counterparty_version,omitempty" yaml:"counterparty_version"`
//	ProofTry              []byte       `protobuf:"bytes,5,opt,name=proof_try,json=proofTry,proto3" json:"proof_try,omitempty" yaml:"proof_try"`
//	ProofHeight           clienttypes.Height `protobuf:"bytes,6,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	Signer                string       `protobuf:"bytes,7,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//
//// MsgChannelOpenAckResponse defines the Msg/ChannelOpenAck response type.
//type MsgChannelOpenAckResponse struct {
//}
//
//
//// MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
//// acknowledge the change of channel state to OPEN on Chain A.
//type MsgChannelOpenConfirm struct {
//	PortId      string       `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
//	ChannelId   string       `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty" yaml:"channel_id"`
//	ProofAck    []byte       `protobuf:"bytes,3,opt,name=proof_ack,json=proofAck,proto3" json:"proof_ack,omitempty" yaml:"proof_ack"`
//	ProofHeight clienttypes.Height `protobuf:"bytes,4,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	Signer      string       `protobuf:"bytes,5,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//// MsgChannelOpenConfirmResponse defines the Msg/ChannelOpenConfirm response type.
//type MsgChannelOpenConfirmResponse struct {
//}
//
//
//// MsgAcknowledgement receives incoming IBC acknowledgement
//type MsgAcknowledgement struct {
//	Packet          Packet       `protobuf:"bytes,1,opt,name=packet,proto3" json:"packet"`
//	Acknowledgement []byte       `protobuf:"bytes,2,opt,name=acknowledgement,proto3" json:"acknowledgement,omitempty"`
//	ProofAcked      []byte       `protobuf:"bytes,3,opt,name=proof_acked,json=proofAcked,proto3" json:"proof_acked,omitempty" yaml:"proof_acked"`
//	ProofHeight     clienttypes.Height `protobuf:"bytes,4,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	Signer          string       `protobuf:"bytes,5,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//// MsgAcknowledgementResponse defines the Msg/Acknowledgement response type.
//type MsgAcknowledgementResponse struct {
//}
//
//// MsgRecvPacket receives incoming IBC packet
//type MsgRecvPacket struct {
//	Packet          Packet       `protobuf:"bytes,1,opt,name=packet,proto3" json:"packet"`
//	ProofCommitment []byte       `protobuf:"bytes,2,opt,name=proof_commitment,json=proofCommitment,proto3" json:"proof_commitment,omitempty" yaml:"proof_commitment"`
//	ProofHeight     clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	Signer          string       `protobuf:"bytes,4,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//// MsgRecvPacketResponse defines the Msg/RecvPacket response type.
//type MsgRecvPacketResponse struct {
//}
//
//
//// MsgTimeout receives timed-out packet
//type MsgTimeout struct {
//	Packet           Packet       `protobuf:"bytes,1,opt,name=packet,proto3" json:"packet"`
//	ProofUnreceived  []byte       `protobuf:"bytes,2,opt,name=proof_unreceived,json=proofUnreceived,proto3" json:"proof_unreceived,omitempty" yaml:"proof_unreceived"`
//	ProofHeight      clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
//	NextSequenceRecv uint64       `protobuf:"varint,4,opt,name=next_sequence_recv,json=nextSequenceRecv,proto3" json:"next_sequence_recv,omitempty" yaml:"next_sequence_recv"`
//	Signer           string       `protobuf:"bytes,5,opt,name=signer,proto3" json:"signer,omitempty"`
//}
//
//
//// MsgTimeoutResponse defines the Msg/Timeout response type.
//type MsgTimeoutResponse struct {
//}
//
