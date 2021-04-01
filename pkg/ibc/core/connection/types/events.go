package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

// IBC connection events
const (
	AttributeKeyConnectionID             = "connection_id"
	AttributeKeyClientID                 = "client_id"
	AttributeKeyCounterpartyClientID     = "counterparty_client_id"
	AttributeKeyCounterpartyConnectionID = "counterparty_connection_id"
)

var EventTypeConnectionOpenInit    = MsgConnectionOpenInit{}.Type()
var EventTypeConnectionOpenTry     = MsgConnectionOpenTry{}.Type()
var EventTypeConnectionOpenAck     = MsgConnectionOpenAck{}.Type()
var EventTypeConnectionOpenConfirm = MsgConnectionOpenConfirm{}.Type()

var AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)

