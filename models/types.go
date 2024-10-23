package models

type TickerInfo struct {
	Name      string
	MarketCap string
	Country   string
	IpoYear   string
	Sector    string
	Industry  string
}

type TickerData struct {
	TickerInfo
	LastPrice       string
	FairValue       string
	MarginOfSafety  string
	EnterpriseValue string
	Roic            string
	OwnerEarnings   string
	TenCap          string
	GrowthData      GrowthData
}

type GrowthData struct {
	RevenueGrowth10Y    string
	RevenueGrowth5Y     string
	EpsGrowth10Y        string
	EpsGrowth5Y         string
	EbitGrowth10Y       string
	EbitGrowth5Y        string
	EbitdaGrowth10Y     string
	EbitdaGrowth5Y      string
	FcfGrowth10Y        string
	FcfGrowth5Y         string
	DividendGrowth10Y   string
	DividendGrowth5Y    string
	BvGrowth10Y         string
	BvGrowth5Y          string
	StockPriceGrowth10Y string
	StockPriceGrowth5Y  string
}

// Define a generic constraint that limits T to string or models.GrowthData
type ScraperReturnDataType interface {
	string | GrowthData
}
