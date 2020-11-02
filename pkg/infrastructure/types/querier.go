package types



type QueryContentParams struct {
	Key       []byte      `json:"key"`
}

func NewContentParams(key []byte) QueryContentParams {
	return QueryContentParams{Key:key}
}