package models

type CalculatedMetrics struct {
	MarketCap     float64 `json:"market_cap"`
	PERatio       float64 `json:"pe_ratio"`
	PBRatio       float64 `json:"pb_ratio"`
	ROE           float64 `json:"roe"`
	DebtEquity    float64 `json:"debt_equity"`
	OPM           float64 `json:"opm"`
	DividendYield float64 `json:"dividend_yield"`
	ROCE          float64 `json:"roce"`
	HighLow       string  `json:"high_low"`
	BookValue     float64 `json:"book_value"`
	FaceValue     float64 `json:"face_value"`
}
