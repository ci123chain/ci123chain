package collactor


var (
	// Ensure that NaiveStrategy satisfies the Strategy interface
	_ Strategy = &NaiveStrategy{}
)
// NewNaiveStrategy returns the proper configs for the NaiveStrategy
func NewNaiveStrategy() *StrategyCfg {
	return &StrategyCfg{
		Type: (&NaiveStrategy{}).GetType(),
	}
}

// NaiveStrategy is an implementation of Strategy.
type NaiveStrategy struct {
	Ordered      bool
	MaxTxSize    uint64 // maximum permitted size of the msgs in a bundled relay transaction
	MaxMsgLength uint64 // maximum amount of messages in a bundled relay transaction
}

func (nrs *NaiveStrategy) HandleEvents(src, dst *Chain, events map[string][]string) {
	panic("implement me")
}

func (nrs *NaiveStrategy) UnrelayedSequences(src, dst *Chain) (*RelaySequences, error) {
	panic("implement me")
}

func (nrs *NaiveStrategy) UnrelayedAcknowledgements(src, dst *Chain) (*RelaySequences, error) {
	panic("implement me")
}

func (nrs *NaiveStrategy) RelayPackets(src, dst *Chain, sp *RelaySequences) error {
	panic("implement me")
}

func (nrs *NaiveStrategy) RelayAcknowledgements(src, dst *Chain, sp *RelaySequences) error {
	panic("implement me")
}

// GetType implements Strategy
func (nrs *NaiveStrategy) GetType() string {
	return "naive"
}


// RelaySequences represents unrelayed packets on src and dst
type RelaySequences struct {
	Src []uint64 `json:"src"`
	Dst []uint64 `json:"dst"`
}
