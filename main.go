package main

import (
	"fmt"
	"os"

	"github.com/zak01011996/bitcoin-ticker/btc_feed"
	"github.com/zak01011996/bitcoin-ticker/cur_feed"
	"github.com/zak01011996/bitcoin-ticker/ticker"
)

// Error log file name
const LOG_FILE = "error.log"

func main() {

	// Prepare bitcoin feeds, add as many as you want...
	btcFeeds := []btcfeed.BtcFeed{
		btcfeed.NewBCFeed("https://blockchain.info/"),
		btcfeed.NewCDFeed("http://api.coindesk.com/"),
	}

	// Prepare currency feeds, add as many as you want...
	currFeeds := []currfeed.CurrFeed{
		currfeed.NewMCNFeed("http://www.mycurrency.net/"),
		currfeed.NewHerokuFeed("http://rate-exchange.herokuapp.com/"),
		currfeed.NewECBFeed("http://www.ecb.europa.eu/"),
	}

	// Prepare error channel for ticker
	errChan := make(chan error)

	// Initialize ticker
	t := &ticker.Ticker{
		BtcFeeds:  btcFeeds,
		CurrFeeds: currFeeds,
		ErrChan:   errChan,
	}

	// Write all errors into file
	go func() {
		// Open file
		f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("ERROR: Cannot open log file: %s\n", err)
			return
		}

		defer f.Close()

		// Listen for error channel and write info into file
		for err := range errChan {
			_, err := f.WriteString(err.Error() + "\n")
			if err != nil {
				fmt.Printf("ERROR: Cannot write to log file: %s\n", err)
			}
		}
	}()

	// Start ticker
	if err := t.Start(); err != nil {
		panic(err)
	}

	// Get result
	t.Print()

	close(errChan)
}
