package ticker

import (
	"fmt"
	"sync"

	"github.com/zak01011996/bitcoin-ticker/btc_feed"
	"github.com/zak01011996/bitcoin-ticker/cur_feed"
)

// Bitcoin realtime ticker
type Ticker struct {
	BtcFeeds  []btcfeed.BtcFeed   // Bitcoin feeds
	CurrFeeds []currfeed.CurrFeed // Currency feeds
	ErrChan   chan error          // Error channel, all errors you can carry here

	wg            sync.WaitGroup // Wait group, to control our swarm
	start         chan struct{}  // To send start action for all routines in one time
	stop          chan struct{}  // To stop response handlers
	btcRes        chan feedResp  // Bitcoin feed result channel
	currRes       chan feedResp  // Currency feed result channel
	totalBtcResp  []feedResp     // To store all responses from bitcoin services
	totalCurrResp []feedResp     // To stora all responses from currency services
}

// Feed response
type feedResp struct {
	Rate float64
	Name string
}

// Will start out ticker
func (t *Ticker) Start() error {

	// Prepare channels
	t.start = make(chan struct{})
	t.stop = make(chan struct{})
	t.btcRes = make(chan feedResp)
	t.currRes = make(chan feedResp)

	// Check if we have feeds on both sides, otherwise ticker is useless
	if len(t.BtcFeeds) == 0 || len(t.CurrFeeds) == 0 {
		return fmt.Errorf("No feeds found: bitcoin feeds(%d), currency feeds(%d)",
			len(t.BtcFeeds), len(t.CurrFeeds))
	}

	totalFeeds := len(t.BtcFeeds) + len(t.CurrFeeds)

	// Add total amount of feeds
	t.wg.Add(totalFeeds)

	// Prepare for race all bitcoin feeds
	for _, f := range t.BtcFeeds {
		go t.processBtc(f)
	}

	// Prepare for race all currency feeds
	for _, f := range t.CurrFeeds {
		go t.processCurr(f)
	}

	// Handle and store responses from bitcoin services
	go func() {
		for {
			select {
			case r := <-t.btcRes:
				t.totalBtcResp = append(t.totalBtcResp, r)
			case <-t.stop:
				return

			}
		}
	}()

	// Handle and store responses from bitcoin services
	go func() {
		for {
			select {
			case r := <-t.currRes:
				t.totalCurrResp = append(t.totalCurrResp, r)
			case <-t.stop:
				return

			}
		}
	}()

	// Try to start all routines in "one time" to make our race more fair...
	for i := 0; i < totalFeeds; i++ {
		t.start <- struct{}{}
	}

	return nil
}

// Prints result
func (t *Ticker) Print() {
	// Wait while all processes will be done
	t.wg.Wait()

	// Try to stop response handlers
	for i := 0; i < 2; i++ {
		t.stop <- struct{}{}
	}

	// Close all channels, cause we've done our JOB
	close(t.start)
	close(t.stop)
	close(t.currRes)
	close(t.btcRes)

	currRespCount := len(t.totalBtcResp)
	btcRespCount := len(t.totalBtcResp)

	// Check response before calculations
	if currRespCount == 0 || btcRespCount == 0 {
		fmt.Printf("No data found from feeds, check error log. bitcoin feeds(%d), currency feeds(%d)\n", btcRespCount, currRespCount)
		return
	}

	// We will calculate using first response, that we got
	btc := t.totalBtcResp[0]
	curr := t.totalCurrResp[0]
	btsUsd := btc.Rate
	eurUsd := curr.Rate
	btcEur := btc.Rate / curr.Rate

	// Print result
	fmt.Printf(
		"BTC/USD: %.2f EUR/USD: %.2f BTC/EUR: %.2f Active sources: BTC/USD (%d of %d, used: %s) EUR/USD (%d of %d, used: %s)\n",
		btsUsd,
		eurUsd,
		btcEur,
		btcRespCount,
		len(t.BtcFeeds),
		btc.Name,
		currRespCount,
		len(t.CurrFeeds),
		curr.Name,
	)
}

func (t *Ticker) processBtc(feed btcfeed.BtcFeed) {
	// Wait for start command
	<-t.start

	// Get result rate from service
	res, err := feed.GetRate(currfeed.USD_CODE)
	if err != nil {
		t.ErrChan <- err
	} else {
		t.btcRes <- feedResp{res, feed.GetName()}
	}

	t.wg.Done()
}

func (t *Ticker) processCurr(feed currfeed.CurrFeed) {
	// Wait for start command
	<-t.start

	// Get result rate from service
	res, err := feed.GetRate(currfeed.EUR_CODE, currfeed.USD_CODE)
	if err != nil {
		t.ErrChan <- err
	} else {
		t.currRes <- feedResp{res, feed.GetName()}
	}

	t.wg.Done()
}
