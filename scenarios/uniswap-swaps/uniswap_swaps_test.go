package uniswapswaps

import (
	"testing"
)

func TestPerTradeSlippage(t *testing.T) {
	// fixed slippage when no band is configured
	s := &Scenario{options: ScenarioOptions{Slippage: 50}}
	if got := s.perTradeSlippage(); got != 50 {
		t.Fatalf("perTradeSlippage = %d, want 50", got)
	}

	// band draws stay within [min, max]
	s.options.SlippageMin = 100
	s.options.SlippageMax = 200
	for i := 0; i < 1000; i++ {
		got := s.perTradeSlippage()
		if got < 100 || got > 200 {
			t.Fatalf("perTradeSlippage = %d, outside [100, 200]", got)
		}
	}

	// tolerances above 100% are capped to keep the output floor non-negative
	s.options.SlippageMin = 0
	s.options.SlippageMax = 0
	s.options.Slippage = 25000
	if got := s.perTradeSlippage(); got != 10000 {
		t.Fatalf("perTradeSlippage = %d, want 10000", got)
	}
}
