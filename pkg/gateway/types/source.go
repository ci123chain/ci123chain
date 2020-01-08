package types

type ServerSource interface {
	FetchSource() []string
}