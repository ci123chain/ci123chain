package types

import (
	"fmt"
	"math/big"
)

type Dec struct {
	i *big.Int
}

// number of decimal places
const (
	Precision = 18

	// bytes required to represent the above precision
	// Ceiling[Log2[999 999 999 999 999 999]]
	DecimalPrecisionBits = 60
)

var (
	precisionReuse       = new(big.Int).Exp(big.NewInt(10), big.NewInt(Precision), nil)
	fivePrecision        = new(big.Int).Quo(precisionReuse, big.NewInt(2))
	precisionMultipliers []*big.Int
	zeroInt              = big.NewInt(0)
	oneInt               = big.NewInt(1)
	tenInt               = big.NewInt(10)
)

func precisionInt() *big.Int {
	return new(big.Int).Set(precisionReuse)
}

func ZeroDec() Dec     { return Dec{new(big.Int).Set(zeroInt)} }
func OneDec() Dec      { return Dec{precisionInt()} }
func SmallestDec() Dec { return Dec{new(big.Int).Set(oneInt)} }

// get the precision multiplier, do not mutate result
func precisionMultiplier(prec int64) *big.Int {
	if prec > Precision {
		panic(fmt.Sprintf("too much precision, maximum %v, provided %v", Precision, prec))
	}
	return precisionMultipliers[prec]
}

// create a new Dec from integer assuming whole number
func NewDec(i int64) Dec {
	return NewDecWithPrec(i, 0)
}

// create a new Dec from integer with decimal place at prec
// CONTRACT: prec <= Precision
func NewDecWithPrec(i, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(big.NewInt(i), precisionMultiplier(prec)),
	}
}

//______________________________________________________________________________________________
//nolint
func (d Dec) IsNil() bool       { return d.i == nil }                 // is decimal nil
func (d Dec) IsZero() bool      { return (d.i).Sign() == 0 }          // is equal to zero
func (d Dec) IsNegative() bool  { return (d.i).Sign() == -1 }         // is negative
func (d Dec) IsPositive() bool  { return (d.i).Sign() == 1 }          // is positive
func (d Dec) Equal(d2 Dec) bool { return (d.i).Cmp(d2.i) == 0 }       // equal decimals
func (d Dec) GT(d2 Dec) bool    { return (d.i).Cmp(d2.i) > 0 }        // greater than
func (d Dec) GTE(d2 Dec) bool   { return (d.i).Cmp(d2.i) >= 0 }       // greater than or equal
func (d Dec) LT(d2 Dec) bool    { return (d.i).Cmp(d2.i) < 0 }        // less than
func (d Dec) LTE(d2 Dec) bool   { return (d.i).Cmp(d2.i) <= 0 }       // less than or equal
func (d Dec) Neg() Dec          { return Dec{new(big.Int).Neg(d.i)} } // reverse the decimal sign
func (d Dec) Abs() Dec          { return Dec{new(big.Int).Abs(d.i)} } // absolute value

// BigInt returns a copy of the underlying big.Int.
func (d Dec) BigInt() *big.Int {
	copy := new(big.Int)
	return copy.Set(d.i)
}

// addition
func (d Dec) Add(d2 Dec) Dec {
	res := new(big.Int).Add(d.i, d2.i)

	if res.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{res}
}

// subtraction
func (d Dec) Sub(d2 Dec) Dec {
	res := new(big.Int).Sub(d.i, d2.i)

	if res.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{res}
}

// multiplication
func (d Dec) MulInt(i Int) Dec {
	mul := new(big.Int).Mul(d.i, i.i)

	if mul.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{mul}
}

// MulInt64 - multiplication with int64
func (d Dec) MulInt64(i int64) Dec {
	mul := new(big.Int).Mul(d.i, big.NewInt(i))

	if mul.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{mul}
}

// quotient
func (d Dec) Quo(d2 Dec) Dec {

	// multiply precision twice
	mul := new(big.Int).Mul(d.i, precisionReuse)
	mul.Mul(mul, precisionReuse)

	quo := new(big.Int).Quo(mul, d2.i)
	chopped := chopPrecisionAndRound(quo)

	if chopped.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{chopped}
}

// quotient truncate
func (d Dec) QuoTruncate(d2 Dec) Dec {

	// multiply precision twice
	mul := new(big.Int).Mul(d.i, precisionReuse)
	mul.Mul(mul, precisionReuse)

	quo := new(big.Int).Quo(mul, d2.i)
	chopped := chopPrecisionAndTruncate(quo)

	if chopped.BitLen() > 255+DecimalPrecisionBits {
		panic("Int overflow")
	}
	return Dec{chopped}
}


// quotient
func (d Dec) QuoInt(i Int) Dec {
	mul := new(big.Int).Quo(d.i, i.i)
	return Dec{mul}
}

// QuoInt64 - quotient with int64
func (d Dec) QuoInt64(i int64) Dec {
	mul := new(big.Int).Quo(d.i, big.NewInt(i))
	return Dec{mul}
}

// TruncateInt truncates the decimals from the number and returns an Int
func (d Dec) TruncateInt() Int {
	return NewIntFromBigInt(chopPrecisionAndTruncateNonMutative(d.i))
}

func chopPrecisionAndTruncateNonMutative(d *big.Int) *big.Int {
	tmp := new(big.Int).Set(d)
	return chopPrecisionAndTruncate(tmp)
}

func chopPrecisionAndRound(d *big.Int) *big.Int {

	// remove the negative and add it back when returning
	if d.Sign() == -1 {
		// make d positive, compute chopped value, and then un-mutate d
		d = d.Neg(d)
		d = chopPrecisionAndRound(d)
		d = d.Neg(d)
		return d
	}

	// get the truncated quotient and remainder
	quo, rem := d, big.NewInt(0)
	quo, rem = quo.QuoRem(d, precisionReuse, rem)

	if rem.Sign() == 0 { // remainder is zero
		return quo
	}

	switch rem.Cmp(fivePrecision) {
	case -1:
		return quo
	case 1:
		return quo.Add(quo, oneInt)
	default: // bankers rounding must take place
		// always round to an even number
		if quo.Bit(0) == 0 {
			return quo
		}
		return quo.Add(quo, oneInt)
	}
}

// create a new Dec from big integer assuming whole numbers
// CONTRACT: prec <= Precision
func NewDecFromInt(i Int) Dec {
	return NewDecFromIntWithPrec(i, 0)
}

// create a new Dec from big integer with decimal place at prec
// CONTRACT: prec <= Precision
func NewDecFromIntWithPrec(i Int, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(i.BigInt(), precisionMultiplier(prec)),
	}
}

// similar to chopPrecisionAndRound, but always rounds down
func chopPrecisionAndTruncate(d *big.Int) *big.Int {
	return d.Quo(d, precisionReuse)
}