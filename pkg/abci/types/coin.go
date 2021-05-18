package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)



func NewChainCoin(amount Int) Coin {
	if err := validate(ChainCoinDenom, amount); err != nil {
		panic(err)
	}

	return Coin{
		Denom:  ChainCoinDenom,
		Amount: amount,
	}
}

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative.
func NewCoin(denom string, amount Int) Coin {

	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

func NewEmptyCoin() Coin {
	res := Coin{
		Denom:  ChainCoinDenom,
		Amount: NewInt(0),
	}
	return res
}

func NewUInt64Coin(denom string, amount uint64) Coin {
	return NewCoin(denom, NewInt(int64(amount)))
}

func (c Coin) String() string {
	by, _ := json.Marshal(c)
	return string(by)
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

func (c Coin) AmountOf(denom string) Int {
	if c.Denom == "" {
		return ZeroInt()
	}

	if c.Denom != denom {
		panic(fmt.Errorf("denom does not match"))
	}
	return c.Amount
}

func (c Coin) Add(coinB Coin) Coin {
	if c.Denom == "" {
		return coinB
	}
	if c.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", c.Denom, coinB.Denom))
	}

	return Coin{coinB.Denom,c.Amount.Add(coinB.Amount)}
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

//-----------------------------------------------------------------------------
// Coins

func init() {
	SetCoinDenomRegex(DefaultCoinDenomRegex)
}

// coinDenomRegex returns the current regex string and can be overwritten for custom validation
var coinDenomRegex = DefaultCoinDenomRegex

// DefaultCoinDenomRegex returns the default regex string
func DefaultCoinDenomRegex() string {
	return reDnmString
}

// SetCoinDenomRegex allows for coin's custom validation by overriding the regular
// expression string used for denom validation.
func SetCoinDenomRegex(reFn func() string) {
	coinDenomRegex = reFn

	reDnm = regexp.MustCompile(fmt.Sprintf(`^%s$`, coinDenomRegex()))
	reDecCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, coinDenomRegex()))
}

// Coins is a set of Coin, one per currency
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

// Empty returns true if there are no coins and false otherwise.
func (coins Coins) Empty() bool {
	return len(coins) == 0
}


// NewCoins constructs a new coin set. The provided coins will be sanitized by removing
// zero coins and sorting the coin set. A panic will occur if the coin set is not valid.
func NewCoins(coins ...Coin) Coins {
	// remove zeroes
	newCoins := removeZeroCoins(coins)
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

func mustValidateDenom(denom string) {
	if err := ValidateDenom(denom); err != nil {
		panic(err)
	}
}


// Returns the amount of a denom from coins
func (coins Coins) AmountOf(denom string) Int {
	mustValidateDenom(denom)

	switch len(coins) {
	case 0:
		return ZeroInt()

	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return ZeroInt()

	default:
		midIdx := len(coins) / 2 // 2:1, 3:1, 4:2
		coin := coins[midIdx]
		switch {
		case denom < coin.Denom:
			return coins[:midIdx].AmountOf(denom)
		case denom == coin.Denom:
			return coin.Amount
		default:
			return coins[midIdx+1:].AmountOf(denom)
		}
	}
}

func (coins Coins) Sub(coinsB Coins) Coins {
	diff, hasNeg := coins.SafeSub(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return diff
}

// SafeSub performs the same arithmetic as Sub but returns a boolean if any
// negative coin amount was returned.
func (coins Coins) SafeSub(coinsB Coins) (Coins, bool) {
	diff := coins.safeAdd(coinsB.negative())
	return diff, diff.IsAnyNegative()
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

func (coins Coins) IsAnyNegative() bool {
	for _, coin := range coins {
		if coin.IsNegative() {
			return true
		}
	}

	return false
}

func (coins Coins) negative() Coins {
	res := make([]Coin, 0, len(coins))

	for _, coin := range coins {
		res = append(res, Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Neg(),
		})
	}

	return res
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

func (coins Coins) Add(coinsB Coins) Coins {
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



func ValidateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}

// validate returns an error if the Coin has a negative amount or if
// the denom is invalid.
func validate(denom string, amount Int) error {
	if err := ValidateDenom(denom); err != nil {
		return err
	}

	if amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", amount)
	}

	return nil
}

// ParseCoinNormalized parses and normalize a cli input for one coin type, returning errors if invalid or on an empty string
// as well.
// Expected format: "{amount}{denomination}"
func ParseCoinNormalized(coinStr string) (coin Coin, err error) {
	decCoin, err := ParseDecCoin(coinStr)
	if err != nil {
		return Coin{}, err
	}

	coin, _ = decCoin.TruncateDecimal()
	return coin, nil
}