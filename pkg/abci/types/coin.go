package types

import "fmt"

type Coin struct {
	Amount Int `json:"amount"`
}

func NewCoin(amount Int) Coin {

	if amount.LT(ZeroInt()) {
		panic(fmt.Errorf("negative coin amount: %v", amount))
	}

	return Coin{
		Amount: amount,
	}
}

func NewUInt64Coin(amount uint64) Coin {
	return NewCoin(NewInt(int64(amount)))
}

func (c Coin) String() string {
	return fmt.Sprintf("%v", c.Amount)
}

func (c Coin) IsZero() bool {
	return c.Amount.IsZero()
}

func (c Coin) IsGTE(other Coin) bool {

	return !c.Amount.LT(other.Amount)
}

func (c Coin) IsLT(other Coin) bool {

	return c.Amount.LT(other.Amount)
}

func (c Coin) IsEqual(other Coin) bool {

	return c.Amount.Equal(other.Amount)
}

func (c Coin) Add(coinB Coin) Coin {

	return Coin{c.Amount.Add(coinB.Amount)}
}

func (c Coin) Sub(coinB Coin) Coin {
	res := Coin{c.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		panic("negative count amount")
	}
	return res
}

func (c Coin) SafeSub(coinB Coin) (Coin, bool) {
	res := Coin{c.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		return res, false
	}
	return res, true
}

func (c Coin) IsNegative() bool {
	return c.Amount.Sign() == -1
}

func (c Coin) IsPositive() bool {
	return c.Amount.Sign() == 1
}

func (c Coin) IsValid() bool {
	return !c.Amount.LT(NewInt(0))
}