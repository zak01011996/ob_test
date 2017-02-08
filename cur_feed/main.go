package currfeed

// Default euro code for currency feed services
const EUR_CODE = "EUR"
const USD_CODE = "USD"

// General currency service interface
type CurrFeed interface {
	GetName() string                         // Name of service
	GetRate(string, string) (float64, error) // Result rate
}
