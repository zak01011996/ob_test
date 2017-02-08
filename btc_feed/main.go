package btcfeed

// General bitcoin service interface
type BtcFeed interface {
	GetName() string                 // Name of service
	GetRate(string) (float64, error) // Result rate in USD
}
