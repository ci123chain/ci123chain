package collactor

import "fmt"

// Strategy defines
type Strategy interface {
	GetType() string
	HandleEvents(src, dst *Chain, events map[string][]string)
	UnrelayedSequences(src, dst *Chain) (*RelaySequences, error)
	UnrelayedAcknowledgements(src, dst *Chain) (*RelaySequences, error)
	RelayPackets(src, dst *Chain, sp *RelaySequences) error
	RelayAcknowledgements(src, dst *Chain, sp *RelaySequences) error
}
// StrategyCfg defines which relaying strategy to take for a given path
type StrategyCfg struct {
	Type string `json:"type" yaml:"type"`
}

// MustGetStrategy returns the strategy and panics on error
func (p *Path) MustGetStrategy() Strategy {
	strategy, err := p.GetStrategy()
	if err != nil {
		panic(err)
	}

	return strategy
}

// GetStrategy the strategy defined in the relay messages
func (p *Path) GetStrategy() (Strategy, error) {
	switch p.Strategy.Type {
	case (&NaiveStrategy{}).GetType():
		return &NaiveStrategy{}, nil
	default:
		return nil, fmt.Errorf("invalid strategy: %s", p.Strategy.Type)
	}
}