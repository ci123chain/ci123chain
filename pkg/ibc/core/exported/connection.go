package exported

type ConnectionI interface {
	GetClientID() string
	GetState() int32
	GetCounterparty() CounterpartyConnectionI
	GetVersions() []Version
	GetDelayPeriod() uint64
	ValidateBasic() error
}

type CounterpartyConnectionI interface {
	GetClientID() string
	GetConnectionID() string
	GetPrefix() Prefix
	ValidateBasic() error
}

type Version interface {
	GetIdentifier() string
	GetFeatures() []string
	VerifyProposedVersion(Version) error
}