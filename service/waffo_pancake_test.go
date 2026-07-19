package service

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/require"
)

func TestNormalizeWaffoPancakeCurrency(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expected   string
		shouldFail bool
	}{
		{name: "blank preserves USD default", input: "", expected: "USD"},
		{name: "lowercase is normalized", input: " cny ", expected: "CNY"},
		{name: "rejects short code", input: "US", shouldFail: true},
		{name: "rejects non letters", input: "US1", shouldFail: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := NormalizeWaffoPancakeCurrency(testCase.input)
			if testCase.shouldFail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.expected, actual)
		})
	}
}

func TestWaffoPancakeTopUpCurrencyFromTradeNo(t *testing.T) {
	require.Equal(t, "CNY", WaffoPancakeTopUpCurrencyFromTradeNo("WAFFO_PANCAKE-CNY-1-2-abc"))
	require.Equal(t, "USD", WaffoPancakeTopUpCurrencyFromTradeNo("WAFFO_PANCAKE-1-2-abc"))
	require.Equal(t, "USD", WaffoPancakeTopUpCurrencyFromTradeNo("WAFFO_PANCAKE_SUB-1-2-abc"))
}

func TestValidateWaffoPancakeTopUpPayment(t *testing.T) {
	testCases := []struct {
		name       string
		tradeNo    string
		money      float64
		currency   string
		amount     string
		shouldFail bool
	}{
		{
			name:     "CNY payment matches encoded order currency and amount",
			tradeNo:  "WAFFO_PANCAKE-CNY-1-2-abc",
			money:    100,
			currency: "CNY",
			amount:   "100.00",
		},
		{
			name:     "legacy order remains USD",
			tradeNo:  "WAFFO_PANCAKE-1-2-abc",
			money:    10,
			currency: "USD",
			amount:   "10",
		},
		{
			name:       "rejects currency mismatch",
			tradeNo:    "WAFFO_PANCAKE-CNY-1-2-abc",
			money:      100,
			currency:   "USD",
			amount:     "100.00",
			shouldFail: true,
		},
		{
			name:       "rejects amount mismatch",
			tradeNo:    "WAFFO_PANCAKE-CNY-1-2-abc",
			money:      100,
			currency:   "CNY",
			amount:     "99.99",
			shouldFail: true,
		},
		{
			name:       "rejects malformed amount",
			tradeNo:    "WAFFO_PANCAKE-CNY-1-2-abc",
			money:      100,
			currency:   "CNY",
			amount:     "invalid",
			shouldFail: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			topUp := &model.TopUp{TradeNo: testCase.tradeNo, Money: testCase.money}
			event := &WaffoPancakeWebhookEvent{Data: WaffoPancakeWebhookData{
				Currency: testCase.currency,
				Amount:   testCase.amount,
			}}
			err := validateWaffoPancakeTopUpPayment(topUp, event)
			if testCase.shouldFail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
