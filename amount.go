package trader

import (
	"fmt"
	"math"

	"strconv"

	"github.com/shopspring/decimal"
)

// Amount represents an amount in a given currency
type Amount struct {
	Trader   *Trader          `json:"-"`
	Value    *decimal.Decimal `json:"value"`
	Currency *Currency        `json:"currency"`
}

// NewAmount creates a new amount structure from a decimal and a currency
func (t *Trader) NewAmount(d *decimal.Decimal, code string) (*Amount, error) {
	c, err := t.Currencies.Find(code)
	if err != nil {
		return nil, err
	}

	return &Amount{
		Trader:   t,
		Value:    d,
		Currency: c,
	}, nil
}

// NewAmountFromFloat creates a new amount structure from a float and
// a currency
func (t *Trader) NewAmountFromFloat(f float64, c string) (*Amount, error) {
	d := decimal.NewFromFloat(f)
	return t.NewAmount(&d, c)
}

// NewAmountFromString creates a new amount structure from a string
// and a currency. Returns an error if the string could not be parsed
func (t *Trader) NewAmountFromString(s, c string) (*Amount, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return t.NewAmount(&d, c)
}

// BaseCurrencyValue returns the value converted to the base currency
// of the Trader. If the Currency of the Amount is already the base currency,
// the value is returned directly
func (a Amount) BaseCurrencyValue() *decimal.Decimal {
	if a.Currency.Is(a.Trader.BaseCurrency.Code) {
		return a.Value
	}

	v := a.Value.Div(*a.Currency.Value)
	return &v
}

// BaseCurrencyAmount returns a new Amount representing the Amount converted
// to the base currency of the Trader. If the Currency of the Amount is already
// the base currency, the Amount is returned directly
func (a Amount) BaseCurrencyAmount() *Amount {
	if a.Currency.Is(a.Trader.BaseCurrency.Code) {
		return &a
	}

	v := a.Value.Div(*a.Currency.Value)
	r, _ := a.Trader.NewAmount(&v, a.Trader.BaseCurrency.Code)
	return r
}

// ToCurrency converts the Amount to the given Currency. If the given Currency
// is the same as the currency one of the Amount, the Amount is returned
// directly
func (a Amount) ToCurrency(code string) (*Amount, error) {
	if a.Currency.Is(code) {
		return &a, nil
	}

	b := a.BaseCurrencyAmount()
	if b.Currency.Is(code) {
		return b, nil
	}

	c, err := a.Trader.Currencies.Find(code)
	if err != nil {
		return nil, err
	}

	v := b.Value.Mul(*c.Value)
	return a.Trader.NewAmount(&v, code)
}

// Add returns a new Amount corresponding to the sum of a and b. The
// Currency of the returned amount is the same as the Currency of a (structure
// on which Add(b) is called). Both a and b are converted to their base currency
// to perform the operation, using their respective BaseCurrencyValues.
// If their base currency differ, an error is returned
func (a Amount) Add(b *Amount) (*Amount, error) {
	if a.Trader.BaseCurrency.Code != b.Trader.BaseCurrency.Code {
		return nil, fmt.Errorf("The base currency of a and b differ: %s & %s",
			a.Trader.BaseCurrency.Code, b.Trader.BaseCurrency.Code)
	}

	c := a.BaseCurrencyValue().Add(*b.BaseCurrencyValue())
	n, _ := a.Trader.NewAmount(&c, a.Trader.BaseCurrency.Code)
	return n.ToCurrency(a.Currency.Code)
}

// Sub returns a new Amount corresponding to the substraction of b from a. The
// Currency of the returned amount is the same as the Currency of a (structure
// on which Sub(b) is called). Both a and b are converted to their base currency
// to perform the operation, using their respective BaseCurrencyValues.
// If their base currency differ, an error is returned
func (a Amount) Sub(b *Amount) (*Amount, error) {
	if a.Trader.BaseCurrency.Code != b.Trader.BaseCurrency.Code {
		return nil, fmt.Errorf("The base currency of a and b differ: %s & %s",
			a.Trader.BaseCurrency.Code, b.Trader.BaseCurrency.Code)
	}

	c := a.BaseCurrencyValue().Sub(*b.BaseCurrencyValue())
	n, _ := a.Trader.NewAmount(&c, a.Trader.BaseCurrency.Code)
	return n.ToCurrency(a.Currency.Code)
}

// Mul returns a new Amount corresponding to the multiplication of a and b. The
// Currency of the returned amount is the same as the Currency of a (structure
// on which Mul(b) is called). Both a and b are converted to their base currency
// to perform the operation, using their respective BaseCurrencyValues.
// If their base currency differ, an error is returned
func (a Amount) Mul(b *Amount) (*Amount, error) {
	if a.Trader.BaseCurrency.Code != b.Trader.BaseCurrency.Code {
		return nil, fmt.Errorf("The base currency of a and b differ: %s & %s",
			a.Trader.BaseCurrency.Code, b.Trader.BaseCurrency.Code)
	}

	c := a.BaseCurrencyValue().Mul(*b.BaseCurrencyValue())
	n, _ := a.Trader.NewAmount(&c, a.Trader.BaseCurrency.Code)
	return n.ToCurrency(a.Currency.Code)
}

// Div returns a new Amount corresponding to the division of a by b. The
// Currency of the returned amount is the same as the Currency of a (structure
// on which Div(b) is called). Both a and b are converted to their base currency
// to perform the operation, using their respective BaseCurrencyValues.
// If their base currency differ, an error is returned.
// Warning: The division isn't precise, the division precision is 16 decimals
func (a Amount) Div(b *Amount) (*Amount, error) {
	if a.Trader.BaseCurrency.Code != b.Trader.BaseCurrency.Code {
		return nil, fmt.Errorf("The base currency of a and b differ: %s & %s",
			a.Trader.BaseCurrency.Code, b.Trader.BaseCurrency.Code)
	}

	c := a.BaseCurrencyValue().Div(*b.BaseCurrencyValue())
	n, _ := a.Trader.NewAmount(&c, a.Trader.BaseCurrency.Code)
	return n.ToCurrency(a.Currency.Code)
}

// Cmp compares a and b precisely in this order.
// Returns:
//	- = 0 if a is equal to b
//  - < 0 if a is smaller than b
//  - > 0 if a is greater than b
// The comparison is done using both amount's base currencies
func (a *Amount) Cmp(b *Amount) (int, error) {
	if a.Trader.BaseCurrency.Code != b.Trader.BaseCurrency.Code {
		return 0, fmt.Errorf("The base currency of a and b differ: %s & %s",
			a.Trader.BaseCurrency.Code, b.Trader.BaseCurrency.Code)
	}

	c := a.BaseCurrencyValue().Cmp(*b.BaseCurrencyValue())
	return c, nil
}

// Round rounds the value to the nearest places
func (a *Amount) Round(places int32) {
	rounded := a.Value.Round(places)
	a.Value = &rounded
}

func (a *Amount) Int64() (int64, error) {
	mulF := math.Pow10(a.Currency.DecimalPlaces())

	factor, err := a.Trader.NewAmountFromFloat(mulF, a.Currency.Code)
	if err != nil {
		return 0, err
	}

	newAmount, err := a.Mul(factor)
	if err != nil {
		return 0, err
	}

	newAmount.Round(0)

	i64, err := strconv.ParseInt(newAmount.String(0), 10, 64)
	if err != nil {
		return 0, err
	}

	return i64, nil

}

// String returns the amount value with the given number of decimals
func (a Amount) String(decimals int32) string {
	return a.Value.StringFixed(decimals)
}
