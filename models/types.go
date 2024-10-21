package models

import "github.com/vue-dot-com/value_investing_screener/parsers/growthnumbers"

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
	GrowthData      growthnumbers.GrowthData
}
