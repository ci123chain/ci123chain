package types

import "math"
// Gas consumption descriptors.
const (
	GasIterNextCostFlatDesc = "IterNextFlat"
	GasValuePerByteDesc     = "ValuePerByte"
	GasWritePerByteDesc     = "WritePerByte"
	GasReadPerByteDesc      = "ReadPerByte"
	GasWriteCostFlatDesc    = "WriteFlat"
	GasReadCostFlatDesc     = "ReadFlat"
	GasHasDesc              = "Has"
	GasDeleteDesc           = "Delete"
	unit					= 1000
)

var (
	cachedKVGasConfig        = KVGasConfig()
	cachedTransientGasConfig = TransientGasConfig()
)

// Gas measured by the SDK
type Gas = uint64

// ErrorOutOfGas defines an error thrown when an action results in out of gas.
type ErrorOutOfGas struct {
	Descriptor string
}

func (e ErrorOutOfGas) Error() string {
	return e.Descriptor
}

// ErrorGasOverflow defines an error thrown when an action results gas consumption
// unsigned integer overflow.
type ErrorGasOverflow struct {
	Descriptor string
}

// GasMeter interface to track gas consumption
type GasMeter interface {
	GasConsumed() Gas
	ConsumeGas(amount Gas, descriptor string)
}

type basicGasMeter struct {
	limit    Gas
	consumed Gas
}

// NewGasMeter returns a reference to a new basicGasMeter.
func NewGasMeter(limit Gas) GasMeter {
	return &basicGasMeter{
		limit:    limit,
		consumed: 0,
	}
}

func (g *basicGasMeter) GasConsumed() Gas {
	return g.consumed
}
// addUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

var gasPrice float64
func SetGasPrice(gp float64) {
	gasPrice = gp
}

func calculateGas(gas Gas) Gas{
	return Gas(math.Ceil(float64(gas)*gasPrice))
}

func (g *basicGasMeter) ConsumeGas(amount Gas, descriptor string) {
	var overflow bool
	//amount /= unit
	// TODO: Should we set the consumed field after overflow checking?
	g.consumed, overflow = addUint64Overflow(g.consumed, calculateGas(amount))
	if overflow {
		panic(ErrorGasOverflow{descriptor})
	}

	if g.consumed > g.limit {
		panic(ErrorOutOfGas{descriptor})
	}
}

type infiniteGasMeter struct {
	consumed Gas
}

// NewInfiniteGasMeter returns a reference to a new infiniteGasMeter.
func NewInfiniteGasMeter() GasMeter {
	return &infiniteGasMeter{
		consumed: 0,
	}
}

func (g *infiniteGasMeter) GasConsumed() Gas {
	return g.consumed
}

func (g *infiniteGasMeter) ConsumeGas(amount Gas, descriptor string) {
	var overflow bool

	// TODO: Should we set the consumed field after overflow checking?
	g.consumed, overflow = addUint64Overflow(g.consumed, calculateGas(amount))
	if overflow {
		panic(ErrorGasOverflow{descriptor})
	}
}

// GasConfig defines gas cost for each operation on KVStores
type GasConfig struct {
	HasCost          Gas
	DeleteCost       Gas
	ReadCostFlat     Gas
	ReadCostPerByte  Gas
	WriteCostFlat    Gas
	WriteCostPerByte Gas
	ValueCostPerByte Gas
	IterNextCostFlat Gas
}

// KVGasConfig returns a default gas configs for KVStores.
func KVGasConfig() GasConfig {
	return GasConfig{
		HasCost:          10,
		DeleteCost:       10,
		ReadCostFlat:     10,
		ReadCostPerByte:  1,
		WriteCostFlat:    10,
		WriteCostPerByte: 10,
		ValueCostPerByte: 1,
		IterNextCostFlat: 15,
	}
}

// TransientGasConfig returns a default gas configs for TransientStores.
func TransientGasConfig() GasConfig {
	// TODO: define gasconfig for transient stores
	return KVGasConfig()
}
