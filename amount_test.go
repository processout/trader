package trader

import (
	"testing"

	"github.com/shopspring/decimal"
)

func getTrader() *Trader {
	usd := decimal.NewFromFloat(1)
	eur := decimal.NewFromFloat(0.8)
	currencies := Currencies{
		NewCurrency("USD", &usd),
		NewCurrency("EUR", &eur),
	}
	trader, _ := New(currencies, "usd")
	return trader
}
func getTrader2() *Trader {
	usd := decimal.NewFromFloat(1)
	eur := decimal.NewFromFloat(0.8)
	currencies := Currencies{
		NewCurrency("USD", &usd),
		NewCurrency("BAD", &eur),
	}
	trader, _ := New(currencies, "bad")
	return trader
}

func TestNewAmount(t *testing.T) {
	trader := getTrader()
	d, _ := decimal.NewFromString("4.2")

	amount, err := trader.NewAmount(&d, "bad")
	if err == nil {
		t.Error("There should have been an error")
	}

	amount, err = trader.NewAmount(&d, "usd")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}

	amount, err = trader.NewAmount(&d, "eur")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}
}

func TestNewAmountFromFloat(t *testing.T) {
	trader := getTrader()

	amount, err := trader.NewAmountFromFloat(4.2, "bad")
	if err == nil {
		t.Error("There should have been an error")
	}

	amount, err = trader.NewAmountFromFloat(4.2, "usd")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}

	amount, err = trader.NewAmountFromFloat(4.2, "eur")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}
}

func TestNewAmountFromString(t *testing.T) {
	trader := getTrader()

	amount, err := trader.NewAmountFromString("bad", "usd")
	if err == nil {
		t.Error("There should have been an error")
	}

	amount, err = trader.NewAmountFromString("4.2", "bad")
	if err == nil {
		t.Error("There should have been an error")
	}

	amount, err = trader.NewAmountFromString("4.2", "usd")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}

	amount, err = trader.NewAmountFromString("4.2", "eur")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if amount == nil {
		t.Error("The returned amount shouldn't have been nil")
	}
	if amount.Trader != trader {
		t.Error("The trader should be the same pointer")
	}
	if amount.String(3) != "4.200" {
		t.Error("Wrong value set: " + amount.String(3))
	}
}

func TestBaseCurrencyValue(t *testing.T) {
	trader := getTrader()
	amount, _ := trader.NewAmountFromString("1", "usd")

	v := amount.BaseCurrencyValue()
	if v != amount.Value {
		t.Error("The base currency value should have stayed the same")
	}

	amount, _ = trader.NewAmountFromString("1", "eur")
	if amount.BaseCurrencyValue().StringFixed(3) != "1.250" {
		t.Error("The base currency value was wrongly converted")
	}
}

func TestBaseCurrencyAmount(t *testing.T) {
	trader := getTrader()
	amount, _ := trader.NewAmountFromString("1", "usd")

	a := amount.BaseCurrencyAmount()
	if a.Value != amount.Value || a.Currency.Code != amount.Currency.Code {
		t.Error("The base currency amount should have stayed the same")
	}

	amount, _ = trader.NewAmountFromString("1", "eur")
	a = amount.BaseCurrencyAmount()
	if a == amount {
		t.Error("The base currency amount should have changed")
	}
	if a.String(3) != "1.250" {
		t.Error("The base currency amount was wrongly converted")
	}
}

func TestToCurrency(t *testing.T) {
	trader := getTrader()
	amount, _ := trader.NewAmountFromString("2", "usd")

	a, err := amount.ToCurrency("bad")
	if err == nil {
		t.Error("There should have been an error")
	}

	a, err = amount.ToCurrency("usd")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if a.Value != amount.Value {
		t.Error("Amount shouldn't have changed")
	}

	a, err = amount.ToCurrency("eur")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if a.String(3) != "1.600" {
		t.Error("The amount was wrongly converted: " + a.String(3))
	}

	a2, err := a.ToCurrency("usd")
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if a2.String(3) != "2.000" {
		t.Error("The amount was wrongly converted: " + a2.String(3))
	}
}

func TestAdd(t *testing.T) {
	trader := getTrader()
	trader2 := getTrader2()
	amount, _ := trader.NewAmountFromString("2.3", "usd")
	amount2, _ := trader2.NewAmountFromString("3.2", "bad")

	s, err := amount2.Add(amount)
	if err == nil {
		t.Error("There should have been an error")
	}

	amount2, _ = trader.NewAmountFromString("3.2", "usd")
	s, err = amount.Add(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "5.500" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	amount2, _ = amount2.ToCurrency("eur")
	s, err = amount.Add(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "5.500" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	s, err = amount2.Add(amount)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "4.400" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("eur") {
		t.Error("The amount currency was incorrectly set")
	}
}

func TestSub(t *testing.T) {
	trader := getTrader()
	trader2 := getTrader2()
	amount, _ := trader.NewAmountFromString("3.2", "usd")
	amount2, _ := trader2.NewAmountFromString("2.3", "bad")

	s, err := amount2.Sub(amount)
	if err == nil {
		t.Error("There should have been an error")
	}

	amount2, _ = trader.NewAmountFromString("2.3", "usd")
	s, err = amount.Sub(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "0.900" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	amount2, _ = amount2.ToCurrency("eur")
	s, err = amount.Sub(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "0.900" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	s, err = amount2.Sub(amount)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "-0.720" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("eur") {
		t.Error("The amount currency was incorrectly set")
	}
}

func TestMul(t *testing.T) {
	trader := getTrader()
	trader2 := getTrader2()
	amount, _ := trader.NewAmountFromString("2.3", "usd")
	amount2, _ := trader2.NewAmountFromString("3.2", "bad")

	s, err := amount2.Mul(amount)
	if err == nil {
		t.Error("There should have been an error")
	}

	amount2, _ = trader.NewAmountFromString("3.2", "usd")
	s, err = amount.Mul(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "7.360" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	amount2, _ = amount2.ToCurrency("eur")
	s, err = amount.Mul(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(3) != "7.360" {
		t.Error("The amount value was incorrectly computed: " + s.String(3))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	s, err = amount2.Mul(amount)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(4) != "5.8880" {
		t.Error("The amount value was incorrectly computed: " + s.String(4))
	}
	if !s.Currency.Is("eur") {
		t.Error("The amount currency was incorrectly set")
	}
}

func TestDiv(t *testing.T) {
	trader := getTrader()
	trader2 := getTrader2()
	amount, _ := trader.NewAmountFromString("2.3", "usd")
	amount2, _ := trader2.NewAmountFromString("3.2", "bad")

	s, err := amount2.Div(amount)
	if err == nil {
		t.Error("There should have been an error")
	}

	amount2, _ = trader.NewAmountFromString("3.2", "usd")
	s, err = amount.Div(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(6) != "0.718750" {
		t.Error("The amount value was incorrectly computed: " + s.String(6))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	amount2, _ = amount2.ToCurrency("eur")
	s, err = amount.Div(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(6) != "0.718750" {
		t.Error("The amount value was incorrectly computed: " + s.String(6))
	}
	if !s.Currency.Is("usd") {
		t.Error("The amount currency was incorrectly set")
	}

	s, err = amount2.Div(amount)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if s == nil {
		t.Error("A new amount should have been returned")
	}
	if s.String(10) != "1.1130434783" {
		t.Error("The amount value was incorrectly computed: " + s.String(10))
	}
	if !s.Currency.Is("eur") {
		t.Error("The amount currency was incorrectly set")
	}
}

func TestCmp(t *testing.T) {
	trader := getTrader()
	trader2 := getTrader2()
	amount, _ := trader.NewAmountFromString("2.3", "usd")
	amount2, _ := trader2.NewAmountFromString("3.2", "bad")

	_, err := amount.Cmp(amount2)
	if err == nil {
		t.Error("There should have been an error")
	}

	amount2, _ = trader.NewAmountFromString("3.2", "usd")
	r, err := amount.Cmp(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r >= 0 {
		t.Error("Answer should have been negative (2.3 < 3.2)")
	}

	amount2, _ = trader.NewAmountFromString("2.14", "usd")
	r, err = amount.Cmp(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r <= 0 {
		t.Error("Answer should have been positive (2.3 > 2.14)")
	}

	amount2, _ = trader.NewAmountFromString("2.3", "usd")
	r, err = amount.Cmp(amount2)
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 0 {
		t.Error("Answer should have been 0 (2.3 == 2.14)")
	}
}

func TestInt64(t *testing.T) {
	trader := getTrader()
	amount, _ := trader.NewAmountFromString("2.3", "usd")

	r, err := amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 230 {
		t.Error("Wrong conversion")
	}

	amount.Currency.Code = "BHD"
	r, err = amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 2300 {
		t.Error("Wrong conversion")
	}

	amount.Currency.Code = "USD"
	amount, _ = trader.NewAmountFromString("2.34", "usd")
	r, err = amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 234 {
		t.Error("Wrong conversion")
	}
}

func TestUint64(t *testing.T) {
	trader := getTrader()
	amount, _ := trader.NewAmountFromString("2.3", "usd")

	r, err := amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 230 {
		t.Error("Wrong conversion")
	}

	amount.Currency.Code = "BHD"
	r, err = amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 2300 {
		t.Error("Wrong conversion")
	}

	amount.Currency.Code = "USD"
	amount, _ = trader.NewAmountFromString("2.34", "usd")
	r, err = amount.Int64()
	if err != nil {
		t.Error("There shouldn't have been an error")
	}
	if r != 234 {
		t.Error("Wrong conversion")
	}
}
