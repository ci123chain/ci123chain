package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCoin(t *testing.T) {
	
	c := NewCoin(DefaultBondDenom ,NewInt(0))
	require.Equal(t, Coin{Denom: DefaultBondDenom ,Amount:NewInt(0)}, c)
}

func TestNewUInt64Coin(t *testing.T) {

	c := NewUInt64Coin(100)
	require.Equal(t, Coin{Amount:NewInt(100)}, c)
}

func TestCoin_String(t *testing.T) {

	c := NewCoin(DefaultBondDenom ,NewInt(0))
	require.Equal(t,"0", c.String())
}

func TestCoin_Add(t *testing.T) {

	c1 := NewUInt64Coin(100)
	c2 := NewUInt64Coin(200)
	c := c1.Add(c2)
	require.Equal(t, NewUInt64Coin(300), c)
}

func TestCoin_Sub(t *testing.T) {

	c1 := NewUInt64Coin(200)
	c2 := NewUInt64Coin(100)
	c := c1.Sub(c2)
	require.Equal(t, NewUInt64Coin(100), c)

	c3 := NewUInt64Coin(100)
	c4 := NewUInt64Coin(200)
	c = c3.Sub(c4) //panic

	c5 := NewUInt64Coin(100)
	c6 := NewUInt64Coin(100)
	c = c5.Sub(c6)
	require.Equal(t, NewUInt64Coin(0), c)
}

func TestCoin_SafeSub(t *testing.T) {

	c1 := NewUInt64Coin(200)
	c2 := NewUInt64Coin(100)
	c, valid := c1.SafeSub(c2)
	fmt.Println(c)
	require.Equal(t, true, valid)
	require.Equal(t, NewUInt64Coin(100), c)

	c3 := NewUInt64Coin(100)
	c4 := NewUInt64Coin(200)
	c, valid = c3.SafeSub(c4)
	require.Equal(t, false, valid)
	fmt.Println(c)

	c5 := NewUInt64Coin(100)
	c6 := NewUInt64Coin(100)
	c, valid = c5.SafeSub(c6)
	fmt.Println(c)
	require.Equal(t, true, valid)
	require.Equal(t, int64(0), c.Amount.Int64())
}

func TestCoin_IsValid(t *testing.T) {

	c1 := NewUInt64Coin(100)
	require.Equal(t, true, c1.IsValid())

	c2 := Coin{Amount:NewInt(-100)}
	require.Equal(t, false, c2.IsValid())
}