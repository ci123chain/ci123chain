package ethereum_gravity

import (
"golang.org/x/crypto/sha3"
"hash"
"sync"
)

func Digest(sig string) []byte {
	k := acquireKeccak()
	k.Write([]byte(sig))
	dst := k.Sum(nil)[:4]
	releaseKeccak(k)
	return dst
}

func acquireKeccak() hash.Hash {
	return keccakPool.Get().(hash.Hash)
}

func releaseKeccak(k hash.Hash) {
	k.Reset()
	keccakPool.Put(k)
}

var keccakPool = sync.Pool{
	New: func() interface{} {
		return sha3.NewLegacyKeccak256()
	},
}