package types

import "github.com/ci123chain/ci123chain/pkg/ibc/core/host"

// IBC client events
const (
	AttributeKeyClientID        = "client_id"
	AttributeKeySubjectClientID = "subject_client_id"
	AttributeKeyClientType      = "client_type"
	AttributeKeyConsensusHeight = "consensus_height"
)
// IBC client events vars
var (
	EventTypeCreateClient         = "create_client"
	EventTypeUpdateClient         = "update_client"
	EventTypeUpgradeClient        = "upgrade_client"
	EventTypeSubmitMisbehaviour   = "client_misbehaviour"
	EventTypeUpdateClientProposal = "update_client_proposal"

	AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)
)

