package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/scripttoken/script/common"
)

var (
	Zero    *big.Int
	Hundred *big.Int
)

func init() {
	Zero = big.NewInt(0)
	Hundred = big.NewInt(100)
}

type Coins struct {
	SCPTWei *big.Int
	SPAYWei *big.Int
}

type CoinsJSON struct {
	SCPTWei *common.JSONBig `json:"scptwei"`
	SPAYWei *common.JSONBig `json:"spaywei"`
}

func NewCoinsJSON(coin Coins) CoinsJSON {
	return CoinsJSON{
		SCPTWei: (*common.JSONBig)(coin.SCPTWei),
		SPAYWei: (*common.JSONBig)(coin.SPAYWei),
	}
}

func (c CoinsJSON) Coins() Coins {
	return Coins{
		SCPTWei: (*big.Int)(c.SCPTWei),
		SPAYWei: (*big.Int)(c.SPAYWei),
	}
}

func (c Coins) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewCoinsJSON(c))
}

func (c *Coins) UnmarshalJSON(data []byte) error {
	var a CoinsJSON
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*c = a.Coins()
	return nil
}

// NewCoins is a convenient method for creating small amount of coins.
func NewCoins(script int64, spay int64) Coins {
	return Coins{
		SCPTWei: big.NewInt(script),
		SPAYWei: big.NewInt(spay),
	}
}

func (coins Coins) String() string {
	return fmt.Sprintf("%v %v, %v %v", coins.SCPTWei, DenomSCPTWei, coins.SPAYWei, DenomSPAYWei)
}

func (coins Coins) IsValid() bool {
	return coins.IsNonnegative()
}

func (coins Coins) NoNil() Coins {
	script := coins.SCPTWei
	if script == nil {
		script = big.NewInt(0)
	}
	spay := coins.SPAYWei
	if spay == nil {
		spay = big.NewInt(0)
	}

	return Coins{
		SCPTWei: script,
		SPAYWei: spay,
	}
}

// CalculatePercentage function calculates amount of coins for the given the percentage
func (coins Coins) CalculatePercentage(percentage uint) Coins {
	c := coins.NoNil()

	p := big.NewInt(int64(percentage))

	script := new(big.Int)
	script.Mul(c.SCPTWei, p)
	script.Div(script, Hundred)

	spay := new(big.Int)
	spay.Mul(c.SPAYWei, p)
	spay.Div(spay, Hundred)

	return Coins{
		SCPTWei: script,
		SPAYWei: spay,
	}
}

// Currently appends an empty coin ...
func (coinsA Coins) Plus(coinsB Coins) Coins {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()

	script := new(big.Int)
	script.Add(cA.SCPTWei, cB.SCPTWei)

	spay := new(big.Int)
	spay.Add(cA.SPAYWei, cB.SPAYWei)

	return Coins{
		SCPTWei: script,
		SPAYWei: spay,
	}
}

func (coins Coins) Negative() Coins {
	c := coins.NoNil()

	script := new(big.Int)
	script.Neg(c.SCPTWei)

	spay := new(big.Int)
	spay.Neg(c.SPAYWei)

	return Coins{
		SCPTWei: script,
		SPAYWei: spay,
	}
}

func (coinsA Coins) Minus(coinsB Coins) Coins {
	return coinsA.Plus(coinsB.Negative())
}

func (coinsA Coins) IsGTE(coinsB Coins) bool {
	diff := coinsA.Minus(coinsB)
	return diff.IsNonnegative()
}

func (coins Coins) IsZero() bool {
	c := coins.NoNil()
	return c.SCPTWei.Cmp(Zero) == 0 && c.SPAYWei.Cmp(Zero) == 0
}

func (coinsA Coins) IsEqual(coinsB Coins) bool {
	cA := coinsA.NoNil()
	cB := coinsB.NoNil()
	return cA.SCPTWei.Cmp(cB.SCPTWei) == 0 && cA.SPAYWei.Cmp(cB.SPAYWei) == 0
}

func (coins Coins) IsPositive() bool {
	c := coins.NoNil()
	return (c.SCPTWei.Cmp(Zero) > 0 && c.SPAYWei.Cmp(Zero) >= 0) ||
		(c.SCPTWei.Cmp(Zero) >= 0 && c.SPAYWei.Cmp(Zero) > 0)
}

func (coins Coins) IsNonnegative() bool {
	c := coins.NoNil()
	return c.SCPTWei.Cmp(Zero) >= 0 && c.SPAYWei.Cmp(Zero) >= 0
}

// ParseCoinAmount parses a string representation of coin amount.
func ParseCoinAmount(in string) (*big.Int, bool) {
	inWei := false
	if len(in) > 3 && strings.EqualFold("wei", in[len(in)-3:]) {
		inWei = true
		in = in[:len(in)-3]
	}

	f, ok := new(big.Float).SetPrec(1024).SetString(in)
	if !ok || f.Sign() < 0 {
		return nil, false
	}

	if !inWei {
		f = f.Mul(f, new(big.Float).SetPrec(1024).SetUint64(1e18))
	}

	ret, _ := f.Int(nil)

	return ret, true
}
