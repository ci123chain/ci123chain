package cosmos_gravity

type Contact struct {
	url string
}

func NewContact(url string) Contact {
	return Contact{url: url}
}

type ChainStatusEnum string

const (
	MOVING ChainStatusEnum = "moving"
	SYNCING ChainStatusEnum = "syncing"
	WAITING_TO_START ChainStatusEnum = "waiting_to_start"
)

type ChainStatus struct {
	BlockHeight uint64
	Status ChainStatusEnum
}

func (c Contact) GetChainStatus() (ChainStatus, error) {
	return ChainStatus{
		BlockHeight: 0,
		Status:      "",
	}, nil
}