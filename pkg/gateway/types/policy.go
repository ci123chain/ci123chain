package types

type LBPolicy interface {
	NextPeer(backends []Instance) Instance
}
