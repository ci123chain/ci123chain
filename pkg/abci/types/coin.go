package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Coin struct {
	Denom  string `json:"denom"`
	Amount Int `json:"amount"`
}

func NewCoin(amount Int) Coin {

	if amount.LT(ZeroInt()) {
		panic(fmt.Errorf("negative coin amount: %v", amount))
	}

	return Coin{
		Denom:"ciCoin",
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
	if c.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", c.Denom, coinB.Denom))
	}

	return Coin{c.Denom,c.Amount.Add(coinB.Amount)}
}

func (c Coin) Sub(coinB Coin) Coin {
	if c.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", c.Denom, coinB.Denom))
	}
	res := Coin{c.Denom,c.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		panic("negative count amount")
	}
	return res
}

func (c Coin) SafeSub(coinB Coin) (Coin, bool) {
	if c.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", c.Denom, coinB.Denom))
	}
	res := Coin{c.Denom,c.Amount.Sub(coinB.Amount)}
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

type Coins []Coin

func (coins Coins) GetDenomByIndex(i int) string {
	return coins[i].Denom
}
// IsValid asserts the Coins are sorted, have positive amount,
// and Denom does not contain upper case characters.
func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		if err := ValidateDenom(coins[0].Denom); err != nil {
			return false
		}
		return coins[0].IsPositive()
	default:
		// check single coin case
		if !(Coins{coins[0]}).IsValid() {
			return false
		}

		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if strings.ToLower(coin.Denom) != coin.Denom {
				return false
			}
			if coin.Denom <= lowDenom {
				return false
			}
			if !coin.IsPositive() {
				return false
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}

		return true
	}
}

func NewCoins(coins ...Coin) Coins {
	// remove zeroes
	newCoins := removeZeroCoins(Coins(coins))
	if len(newCoins) == 0 {
		return Coins{}
	}

	newCoins.Sort()

	// detect duplicate Denoms
	if dupIndex := findDup(newCoins); dupIndex != -1 {
		panic(fmt.Errorf("find duplicate denom: %s", newCoins[dupIndex]))
	}

	if !newCoins.IsValid() {
		panic(fmt.Errorf("invalid coin set: %s", newCoins))
	}

	return newCoins
}

// removeZeroCoins removes all zero coins from the given coin set in-place.
func removeZeroCoins(coins Coins) Coins {
	i, l := 0, len(coins)
	for i < l {
		if coins[i].IsZero() {
			// remove coin
			coins = append(coins[:i], coins[i+1:]...)
			l--
		} else {
			i++
		}
	}

	return coins[:i]
}

func (coins Coins) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}
	return out[:len(out)-1]
}

func (coins Coins) Len() int           { return len(coins) }
func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }
func (coins Coins) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

// Sort is a helper function to sort the set of coins inplace
func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}

// safeAdd will perform addition of two coins sets. If both coin sets are
// empty, then an empty set is returned. If only a single set is empty, the
// other set is returned. Otherwise, the coins are compared in order of their
// denomination and addition only occurs when the denominations match, otherwise
// the coin is simply added to the sum assuming it's not zero.
func (coins Coins) safeAdd(coinsB Coins) Coins {
	sum := ([]Coin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				// return nil coins if both sets are empty
				return sum
			}

			// return set B (excluding zero coins) if set A is empty
			return append(sum, removeZeroCoins(coinsB[indexB:])...)
		} else if indexB == lenB {
			// return set A (excluding zero coins) if set B is empty
			return append(sum, removeZeroCoins(coins[indexA:])...)
		}

		coinA, coinB := coins[indexA], coinsB[indexB]

		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1: // coin A denom < coin B denom
			if !coinA.IsZero() {
				sum = append(sum, coinA)
			}

			indexA++

		case 0: // coin A denom == coin B denom
			res := coinA.Add(coinB)
			if !res.IsZero() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1: // coin A denom > coin B denom
			if !coinB.IsZero() {
				sum = append(sum, coinB)
			}

			indexB++
		}
	}
}

// IsZero returns whether all coins are zero
func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.Amount.IsZero() {
			return false
		}
	}
	return true
}

func (coins Coins) Add(coinsB ...Coin) Coins {
	return coins.safeAdd(coinsB)
}

type findDupDescriptor interface {
	GetDenomByIndex(int) string
	Len() int
}

func findDup(coins findDupDescriptor) int {
	if coins.Len() <= 1 {
		return -1
	}

	prevDenom := coins.GetDenomByIndex(0)
	for i := 1; i < coins.Len(); i++ {
		if coins.GetDenomByIndex(i) == prevDenom {
			return i
		}
		prevDenom = coins.GetDenomByIndex(i)
	}

	return -1
}

var (
	reDnmString = `[a-z][a-z0-9/]{2,31}`
	reDnm       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDnmString))
)

func ValidateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}