package types


type ProxyType string



const (
	LBType = "lb"
	ConcretType = "concret"
	FilterType = "filter"
	LB = LBType
	Concret = ConcretType
	Filter = FilterType
	ValidCode uint64 = 200
)

type ResultRep struct {
	Code    uint64     `json:"coee"`
	Data    string     `json:"data"`
}
