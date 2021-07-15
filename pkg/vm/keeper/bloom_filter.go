package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (k *Keeper) RecordSection(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	index := height/evmtypes.SectionSize

	gen, found := k.GetSectionBloom(ctx, index)
	if !found {
		gen, _ = NewGenerator(evmtypes.SectionSize)
	}
	gen.AddBloom(uint(height-index * evmtypes.SectionSize-1), bloom)

	k.SetSectionBloom(ctx, index, gen)
}

func (k *Keeper) GetSectionBloom(ctx sdk.Context, index int64) (*Generator, bool) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixSection)

	has := store.Has(evmtypes.BloomKey(index))
	if !has {
		return nil, false
	}

	bz := store.Get(evmtypes.BloomKey(index))
	var section Generator
	err := json.Unmarshal(bz, &section)
	if err != nil {
		return nil, false
	}
	return &section, true
}

func (k *Keeper) SetSectionBloom(ctx sdk.Context, index int64, gen *Generator) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixSection)
	by, _ := json.Marshal(gen)
	store.Set(evmtypes.BloomKey(index), by)
}

var (
	// errSectionOutOfBounds is returned if the user tried to add more bloom filters
	// to the batch than available space, or if tries to retrieve above the capacity.
	errSectionOutOfBounds = evmtypes.ErrSectionOutOfBounds

	// errBloomBitOutOfBounds is returned if the user tried to retrieve specified
	// bit bloom above the capacity.
	errBloomBitOutOfBounds = evmtypes.ErrBloomBitOutOfBounds
)

// Generator takes a number of bloom filters and generates the rotated bloom bits
// to be used for batched filtering.
type Generator struct {
	Blooms   [ethtypes.BloomBitLength][]byte `json:"blooms"`// Rotated blooms for per-bit matching
	Sections uint  `json:"sections"`                       // Number of sections to batch together
	//NextSec  uint   `json:"next_sec"`                     // Next section to set when adding a bloom
}

// NewGenerator creates a rotated bloom generator that can iteratively fill a
// batched bloom filter's bits.
func NewGenerator(sections uint) (*Generator, error) {
	if sections % 8 != 0 {
		return nil, evmtypes.ErrBloomFilterSectionNum
	}
	b := &Generator{Sections: sections}
	for i := 0; i < ethtypes.BloomBitLength; i++ {
		b.Blooms[i] = make([]byte, sections/8)
	}
	return b, nil
}

// AddBloom takes a single bloom filter and sets the corresponding bit column
// in memory accordingly.
func (b *Generator) AddBloom(index uint, bloom ethtypes.Bloom) error {
	// Make sure we're not adding more bloom filters than our capacity
	if index >= b.Sections {
		return errSectionOutOfBounds
	}
	//if b.NextSec != index {
	//	return errors.New("bloom filter with unexpected index")
	//}
	// Rotate the bloom and insert into our collection
	byteIndex := index / 8
	bitMask := byte(1) << byte(7-index%8)

	for i := 0; i < ethtypes.BloomBitLength; i++ {
		bloomByteIndex := ethtypes.BloomByteLength - 1 - i/8
		bloomBitMask := byte(1) << byte(i%8)

		if (bloom[bloomByteIndex] & bloomBitMask) != 0 {
			b.Blooms[i][byteIndex] |= bitMask
		}
	}
	//b.NextSec++

	return nil
}

// Bitset returns the bit vector belonging to the given bit index after all
// blooms have been added.
func (b *Generator) Bitset(idx uint) ([]byte, error) {
	//if b.NextSec != b.Sections {
	//	return nil, errors.New("bloom not fully generated yet")
	//}
	if idx >= ethtypes.BloomBitLength {
		return nil, errBloomBitOutOfBounds
	}
	return b.Blooms[idx], nil
}

