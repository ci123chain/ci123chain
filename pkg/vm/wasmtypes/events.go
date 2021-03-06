package types

import vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"

const (
	EventType = "contract"

	//EventTypeUpload = "upload_contract"
	//EventTypeInitiate = "init_contract"
	//EventTypeInvoke = "invoke_contract"
	//EventTypeMigrate = "migrate_contract"


	AttributeKeyCodeHash = "code_hash"
	AttributeKeyAddress = "contract_address"
	AttributeKeyMethod = "operation"
	AttributeValueUpload = "upload_contract"
	AttributeValueInitiate = "init_contract"
	AttributeValueInvoke = "invoke_contract"
	AttributeValueMigrate = "migrate_contract"

	AttributeValueCategory        = vmmodule.ModuleName
)
