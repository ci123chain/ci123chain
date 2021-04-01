package store

import (
	"bytes"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Gets the first item.
func First(st KVStore, start, end []byte) (kv abci.EventAttribute, ok bool) {
	iter := st.Iterator(start, end)
	if !iter.Valid() {
		return kv, false
	}
	defer iter.Close()

	return abci.EventAttribute{Key: iter.Key(), Value: iter.Value()}, true
}

// Gets the last item.  `end` is exclusive.
func Last(st KVStore, start, end []byte) (kv abci.EventAttribute, ok bool) {
	iter := st.ReverseIterator(end, start)
	if !iter.Valid() {
		if v := st.Get(start); v != nil {
			return abci.EventAttribute{Key: cp(start), Value: cp(v)}, true
		}
		return kv, false
	}
	defer iter.Close()

	if bytes.Equal(iter.Key(), end) {
		// Skip this one, end is exclusive.
		iter.Next()
		if !iter.Valid() {
			return kv, false
		}
	}

	return abci.EventAttribute{Key: iter.Key(), Value: iter.Value()}, true
}
