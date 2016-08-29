package trader

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Currency represents a currency and its value relative to the dollar
type Currency struct {
	// Code is the ISO 4217 code of the currency
	Code string `json:"code"`
	// Value is the value of the currency, relative to the base currency
	Value *decimal.Decimal `json:"value"`
}

// NewCurrency creates a new Currency structure
func NewCurrency(code string, v *decimal.Decimal) *Currency {
	return &Currency{
		Code:  strings.ToUpper(code),
		Value: v,
	}
}

// Currencies represents a slice of Currencies
type Currencies []*Currency

// Equals compares both currencies to see if they're Equals
func (c Currencies) Equals(c2 Currencies) bool {
	if c == c2 {
		return true
	}
	if c == nil || c2 == nil {
		return false
	}

	for _, v := range c {
		found := false
		for _, v2 := range c2 {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// Find finds a Currency within the Currencies slice from the given
// currency code, or returns an error if the currency code was not found
func (c Currencies) Find(code string) (*Currency, error) {
	for _, v := range c {
		if v.Is(code) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The currency code %s could not be found.", code)
}

// Is returns true if the given code is the code of the Currency, false
// otherwise
func (c Currency) Is(code string) bool {
	return c.Code == strings.ToUpper(code)
}

// DecimalPlaces returns the number of decimal places a currency has
// e.g. for USD there are 2 ($12.25), for JPY there are 0 (5412)
func (c Currency) DecimalPlaces() int {
	// Here we just test for the currencies that don't have 2 decimal places

	switch c.Code {

	case "BIF", "BYR", "CLP", "DJF", "GNF", "ISK", "JPY", "KMF", "KRW",
		"XPF", "XOF", "XAF", "VUV", "VND", "UYI", "UGX", "RWF", "PYG":
		return 0

	case "BHD", "IQD", "JOD", "KWD", "LYD", "TND", "OMR":
		return 3

	case "CLF":
		return 4

	default:
		return 2

	}
}
