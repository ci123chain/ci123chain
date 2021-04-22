package types

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"context"
)

type MsgServer interface {
	// Unjail defines a method for unjailing a jailed validator, thus returning
	// them into the bonded validator set, so they can begin receiving provisions
	// and rewards again.
	Unjail(context.Context, *MsgUnjail) (*MsgUnjailResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) Unjail(ctx context.Context, req *MsgUnjail) (*MsgUnjailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unjail not implemented")
}

type MsgUnjailResponse struct {
}
