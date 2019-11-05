package types

import "fmt"

type Coin uint64

func NewCoin() Coin {
	return 0
}

func (c Coin) String() string {
	return fmt.Sprintf("%f", c)
}

func (c Coin) Sub(coinB Coin) Coin {
	return c - coinB
}

func (c Coin) Add(coinB Coin) Coin {
	return c.SafeAdd(coinB)
}

func (c Coin) SafeSub(coinB Coin) (Coin, bool) {
	res := c - coinB
	return res, res.IsValid()
}

func (c Coin) SafeAdd(coinB Coin) Coin {
	return c + coinB
}

func (c Coin)IsValid() bool {
	return c >= 0
}

func (c Coin) IsAnyNegative() bool {
	return c < 0
}