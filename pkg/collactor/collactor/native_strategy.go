package collactor


var (
	// Ensure that NaiveStrategy satisfies the Strategy interface
	_ Strategy = &NaiveStrategy{}
)
// NewNaiveStrategy returns the proper config for the NaiveStrategy
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

func (nrs *NaiveStrategy) UnrelayedSequences(src, dst *Chain) (*interface{}, error) {
	panic("implement me")
}

func (nrs *NaiveStrategy) UnrelayedAcknowledgements(src, dst *Chain) (*interface{}, error) {
	panic("implement me")
}

func (nrs *NaiveStrategy) RelayPackets(src, dst *Chain, sp *interface{}) error {
	panic("implement me")
}

func (nrs *NaiveStrategy) RelayAcknowledgements(src, dst *Chain, sp *interface{}) error {
	panic("implement me")
}

// GetType implements Strategy
func (nrs *NaiveStrategy) GetType() string {
	return "naive"
}