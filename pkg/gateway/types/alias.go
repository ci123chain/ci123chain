package types


type ProxyType string



const (
	LBType = "lb"
	ConcretType = "concret"
	FilterType = "filter"
	LB = LBType
	Concret = ConcretType
	Filter = FilterType
	ValidCode int = 200

	ErrGetErrorResponse uint32 = 101
	ErrUnmarshalFailed uint32 = 102

)

type ResultRep struct {
	Code    uint64     `json:"coee"`
	Data    string     `json:"data"`
}


type ErrorResponse struct {
	Err string `json:"err"`
	//Code uint32 `json:"code"`
}