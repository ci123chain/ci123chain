package exported

type Root interface {
	GetHash() []byte
	Empty() bool
}

type Prefix interface {
	Bytes() []byte
	Empty() bool
}

// Path implements spec:CommitmentPath.
// A path is the additional information provided to the verification function.
type Path interface {
	String() string
	Empty() bool
}
