package engine

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"coinmarketscraper/coin"

	"golang.org/x/net/html"
)

func Run(cancelChan <-chan struct{}) <-chan coin.Coin {
	cancelChans := []chan struct{}{}
	cc := make(chan struct{})
	ch := newCoinScraper("https://coinmarketcap.com").startScrap(cc)
	cancelChans = append(cancelChans, cc)
	for i := 1; i < 3; i++ {
		cc = make(chan struct{})
		ch = merge(ch, newCoinScraper(fmt.Sprintf("https://coinmarketcap.com/%d", i)).
			startScrap(cc))
		cancelChans = append(cancelChans, cc)
	}
	go func() {
		<-cancelChan
		for _, c := range cancelChans {
			c <- struct{}{}
			close(c)
		}
	}()
	return ch
}

type coinScraper struct {
	url    string
	client http.Client
}

func newCoinScraper(url string) coinScraper {
	return coinScraper{
		url: url,
		client: http.Client{
			Timeout: time.Duration(2 * time.Second),
		},
	}
}

func (cs coinScraper) startScrap(cancelChan <-chan struct{}) <-chan coin.Coin {
	output := make(chan coin.Coin, 10)
	go func() {
		ticker := time.NewTicker(1 * time.Second).C
		for {
			select {
			case <-cancelChan:
				fmt.Println("Canceled scraping: ", cs.url)
				close(output)
				return
			case <-ticker:
				err := cs.scrap(output)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}()

	return output
}

func (cs coinScraper) scrap(output chan<- coin.Coin) error {
	resp, err := cs.client.Get(cs.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			isTbody := t.Data == "tbody"
			if isTbody {
				scrapCoins(z, output)
			}

		}
	}
}

func scrapCoins(token *html.Tokenizer, output chan<- coin.Coin) {
	for {
		tt := token.Next()
		switch {
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == "tr" {
				output <- scrapCoin(token, t.Attr[0].Val)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "tbody" {
				return
			}
		}
	}
}

func scrapCoin(token *html.Tokenizer, id string) coin.Coin {
	coinInfo := [10]string{}
	count := 0
	for {
		tt := token.Next()
		switch {
		case tt == html.TextToken:
			t := token.Token()
			dataClean := strings.Replace(strings.Replace(t.Data, " ", "", -1), "\n", "", -1)
			if dataClean != "" {
				coinInfo[count] = dataClean
				count++
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "tr" {
				return buildCoin(id, coinInfo)
			}
		}
	}
}

func buildCoin(coinID string, coinInfo [10]string) coin.Coin {
	coin := coin.Coin{
		ID:         coinID,
		Name:       coinInfo[2],
		Acronym:    coinInfo[1],
		IsMineable: false,
	}
	if pos, err := strconv.Atoi(coinInfo[0]); err == nil {
		coin.RankPosition = pos
	}
	if mc, err := parseDollar(coinInfo[3]); err == nil {
		coin.MarketCap = mc
	}
	if price, err := parseDollar(coinInfo[4]); err == nil {
		coin.Price = price
	}
	if vol, err := parseDollar(coinInfo[5]); err == nil {
		coin.Volume24h = vol
	}
	if cir, err := parseDollar(coinInfo[6]); err == nil {
		coin.CirculatingSupply = cir
	}
	if coinInfo[8] != "*" {
		coin.IsMineable = true
		if change, err := strconv.ParseFloat(coinInfo[8][0:len(coinInfo[8])-1], 64); err == nil {
			coin.Change24h = change
		}
	} else {
		coin.IsMineable = false
		if change, err := strconv.ParseFloat(coinInfo[9][0:len(coinInfo[9])-1], 64); err == nil {
			coin.Change24h = change
		}
	}

	return coin
}

func parseDollar(dollarValue string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(dollarValue[1:], ",", "", -1), 64)
}

func merge(cs ...<-chan coin.Coin) <-chan coin.Coin {
	var wg sync.WaitGroup
	out := make(chan coin.Coin, 100)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan coin.Coin) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
